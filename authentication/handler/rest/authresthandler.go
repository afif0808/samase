package authenticationresthandler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fifentory/options"
	"net/http"
	authenticationmiddleware "samase/authentication/middleware"
	"samase/jsonwebtoken"
	jsonwebtokenredisrepo "samase/jsonwebtoken/repository/redis"
	samasemodels "samase/models"
	userrepo "samase/user/repository"
	usersqlrepo "samase/user/repository/sql"
	"strings"
	"time"

	"google.golang.org/api/oauth2/v2"

	"github.com/gomodule/redigo/redis"

	authenticationservice "samase/authentication/service"
	jsonwebtokenservice "samase/jsonwebtoken/service"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func InjectAuthenticationRESTHandler(conn *sql.DB, ee *echo.Echo, redisConn redis.Conn) {
	ussf := usersqlrepo.NewUserSQLFetcher(conn)

	createJWT := jsonwebtokenservice.CreateJWT()
	secretKey := []byte("itssignaturekey")
	jwtsm := jwt.SigningMethodHS256
	tokenDuration := time.Hour * 30

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
	ee.GET("/authenticate/android", AndroidAuthenticate(authenticationmiddleware.InjectAuthenticate()))
	blackListJWT := jsonwebtokenredisrepo.BlackListJWT(redisConn)
	logout := authenticationservice.Logout(blackListJWT)
	ee.POST("/logout", Logout(logout))
}

func AndroidAuthenticate(authenticateMiddleware echo.MiddlewareFunc) echo.HandlerFunc {
	return authenticateMiddleware(func(ectx echo.Context) error {
		// ctx := ectx.Request().Context()
		return ectx.JSON(http.StatusOK, ectx.Get("user"))
	})
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

		return ectx.JSON(http.StatusOK, struct {
			Token string `json:"token"`
		}{Token: tokenStr})
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

func verifyIdToken(ctx context.Context, idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.NewService(ctx)
	if err != nil {
		return nil, err
	}
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)

	tokenInfo, err := tokenInfoCall.Do()

	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
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
			IDToken string `json:"id_token"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusUnauthorized, nil)
		}

		var googleUser struct {
			Email      string `json:"email"`
			Registered bool   `json:"registered"`
			Token      string `json:"token,omitempty"`
		}

		tokenInfo, err := verifyIdToken(ctx, post.IDToken)
		if err != nil {
			return ectx.JSON(http.StatusUnauthorized, nil)
		}

		googleUser.Email = tokenInfo.Email

		usfe.WithEmail()
		uss, err := usfe.GetUsers(ctx, &options.Options{
			Filters: []options.Filter{options.Filter{
				By:       "user_email.value",
				Operator: "=",
				Value:    tokenInfo.Email,
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
		authorizationBearer := ectx.Request().Header.Get("Authorization")
		if !strings.Contains(authorizationBearer, "Bearer ") {
			return ectx.JSON(http.StatusUnauthorized, nil)
		}
		token := strings.Replace(authorizationBearer, "Bearer ", "", 1)
		err := logout(token)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, samasemodels.RESTResponse{Message: err})
		}
		return ectx.JSON(http.StatusOK, nil)
	}
}
