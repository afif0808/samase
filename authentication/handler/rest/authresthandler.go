package authenticationresthandler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fifentory/options"
	"net/http"
	"samase/jsonwebtoken"
	jsonwebtokenredisrepo "samase/jsonwebtoken/repository/redis"
	samasemodels "samase/models"
	userrepo "samase/user/repository"
	usersqlrepo "samase/user/repository/sql"
	"time"

	"github.com/gomodule/redigo/redis"

	authenticationservice "samase/authentication/service"
	jsonwebtokenservice "samase/jsonwebtoken/service"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func InjectAuthenticationRESTHandler(conn *sql.DB, ee *echo.Echo, redisConn redis.Conn) {
	ussf := usersqlrepo.NewUserSQLFetcher(conn)

	createJWT := jsonwebtokenservice.CreateJWT()
	secretKey := []byte("itssignaturekey")
	jwtsm := jwt.SigningMethodHS256
	tokenDuration := time.Minute * 30

	ee.POST(
		"/login",
		Login(
			&ussf,
			createJWT,
			secretKey,
			jwtsm,
			tokenDuration,
		),
	)
	ee.POST("/login/google", GoogleLogin(&ussf, createJWT, secretKey, jwtsm, tokenDuration))

	blackListJWT := jsonwebtokenredisrepo.BlackListJWT(redisConn)
	logout := authenticationservice.Logout(blackListJWT)
	ee.POST("/logout", Logout(logout))
}

func Login(
	usfe userrepo.UserFetcher,
	createJWT jsonwebtoken.CreateJWTFunc,
	secretKey interface{},
	jwtsm jwt.SigningMethod,
	tokenDuration time.Duration,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		if ctx == nil {
			ctx = context.Background()
		}
		if _, isWithDeadline := ctx.Deadline(); !isWithDeadline {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Second*5)
			defer cancel()
		}
		var post struct {
			Name     string `json:"name"`
			Password string `json:"password"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusUnauthorized, samasemodels.RESTResponse{Message: "login failed"})
		}
		usfe.WithPassword()
		uss, err := usfe.GetUsers(ctx, &options.Options{
			Filters: []options.Filter{
				options.Filter{
					By:       "user.name",
					Operator: "=",
					Value:    post.Name,
				},
			},
		})
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, samasemodels.RESTResponse{Message: "there's a problem"})
		}
		if len(uss) <= 0 {
			return ectx.JSON(http.StatusUnauthorized, samasemodels.RESTResponse{Message: "login failed ,username or password is incorrect "})
		}

		us := uss[0]
		if !checkPasswordHash(us.Password.Hash, post.Password) {
			return ectx.JSON(http.StatusUnauthorized, samasemodels.RESTResponse{Message: "login failed ,username or password is incorrect "})
		}

		jwtscl := jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
		}
		sajwtcl := jsonwebtoken.SamaseJWTClaims{
			StandardClaims: jwtscl,
			User:           uss[0],
		}

		tokenStr, err := createJWT(sajwtcl, secretKey, jwtsm)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, samasemodels.RESTResponse{Message: "there's a problem"})
		}

		ectx.Response().Header().Set("X-AUTH-TOKEN", tokenStr)
		return ectx.JSON(http.StatusOK, "Login success")
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func checkPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func GoogleLogin(
	usfe userrepo.UserFetcher,
	createJWT jsonwebtoken.CreateJWTFunc,
	secretKey interface{},
	jwtsm jwt.SigningMethod,
	tokenDuration time.Duration,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			AccessToken string `json:"access_token"`
		}
		err := ectx.Bind(&post)

		if err != nil {
			return ectx.JSON(http.StatusUnauthorized, nil)
		}
		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + post.AccessToken)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		var googleUser struct {
			Email      string `json:"email"`
			Registered bool   `json:"registered"`
			Token      string `json:"token,omitempty"`
		}
		err = json.NewDecoder(resp.Body).Decode(&googleUser)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}

		usfe.WithEmail()
		uss, err := usfe.GetUsers(ctx, &options.Options{
			Filters: []options.Filter{options.Filter{
				By:       "user_email.value",
				Operator: "=",
				Value:    googleUser.Email,
			}},
		})
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}

		if len(uss) > 0 {
			googleUser.Registered = true
			jwtcl := jsonwebtoken.SamaseJWTClaims{
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(tokenDuration).Unix(),
				},
				User: uss[0],
			}

			tokenStr, err := createJWT(jwtcl, secretKey, jwtsm)
			googleUser.Token = tokenStr
			if err != nil {
				return ectx.JSON(http.StatusInternalServerError, nil)
			}
		}

		return ectx.JSON(http.StatusOK, googleUser)
	}
}

func Logout(
	logout authenticationservice.LogoutFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		var post struct {
			AccessToken string `json:"access_token"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusBadRequest, samasemodels.RESTResponse{Message: err})
		}
		err = logout(post.AccessToken)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, samasemodels.RESTResponse{Message: err})
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}

type googleAuthStuffs struct {
	dataURL string
	code    string
}

func getGoogleData(ctx context.Context, conf *oauth2.Config, gass googleAuthStuffs) (*http.Response, error) {
	token, err := conf.Exchange(oauth2.NoContext, gass.code)
	if err != nil {
		return nil, err
	}

	cl := conf.Client(ctx, token)
	res, err := cl.Get(gass.dataURL)
	return res, err
}
