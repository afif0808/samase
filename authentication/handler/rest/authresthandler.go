package authenticationresthandler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	authenticationmiddleware "samase/authentication/middleware"
	"samase/jsonwebtoken"
	jsonwebtokenrepo "samase/jsonwebtoken/repository"
	jsonwebtokenredisrepo "samase/jsonwebtoken/repository/redis"
	samasemodels "samase/models"
	samasemailservice "samase/samasemail/service"
	"samase/user"
	usersqlrepo "samase/user/repository/sql"
	userservice "samase/user/service"
	"samase/useremail"
	useremailsqlrepo "samase/useremail/repository/sql"
	userpasswordsqlrepo "samase/userpassword/repository/sql"
	"strings"
	"time"

	"gopkg.in/gomail.v2"

	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"

	authenticationservice "samase/authentication/service"
	jsonwebtokenservice "samase/jsonwebtoken/service"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func InjectAuthenticationRESTHandler(conn *sql.DB, ee *echo.Echo, redisConn redis.Conn) {
	gussf := usersqlrepo.GetUserSQLFetcher(conn)
	// ussf := usersqlrepo.NewUserSQLFetcher(conn)
	secretKey := []byte("itssignaturekey")
	jwtsm := jwt.SigningMethodHS256
	createJWT := jsonwebtokenservice.CreateJWT(secretKey, jwtsm)
	parseJWT := jsonwebtokenservice.ParseJWT(secretKey, jwtsm)
	tokenDuration := time.Hour * 30
	login := authenticationservice.Login(usersqlrepo.GetUserSQLFetcher(conn))
	ee.POST(
		"/login",
		Login(
			login,
			createJWT,
			tokenDuration,
		),
	)
	audience := "744967159273-rhjtp67vu4075un8hftrp5silgbh2n6f.apps.googleusercontent.com"
	credentialFile := "/root/another-gogole.json"
	verifyIDToken := authenticationservice.GoogleVerifyIDToken(audience, credentialFile)

	createUser := usersqlrepo.CreateUser(conn)
	createUserEmail := useremailsqlrepo.CreateUserEmail(conn)
	createUserPassword := userpasswordsqlrepo.CreateUserPassword(conn)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", mailer.FormatAddress("afifsamase@gmail.com", "SamaseApp"))
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "afifsamase@gmail.com", "afifsamase0808")
	sendEmail := samasemailservice.SendEmail(dialer, mailer)
	ee.POST("/login/google", GoogleLogin(userservice.CreateUser(createUser, createUserEmail, createUserPassword), userservice.GetUserByEmail(gussf), verifyIDToken, createJWT, tokenDuration, sendEmail))
	ee.GET("/authenticate/android", AndroidAuthenticate(authenticationmiddleware.InjectAuthenticate()))
	blackListJWT := jsonwebtokenredisrepo.BlackListJWT(redisConn)
	logout := authenticationservice.Logout(blackListJWT)
	ee.GET("/logout", Logout(logout))

	getUserByID := userservice.GetUserByID(gussf)

	refreshToken := authenticationservice.RefreshToken(createJWT, parseJWT, getUserByID)

	ee.GET("/token/refresh", RefreshToken(refreshToken, time.Hour*20000, blackListJWT))
}

func AndroidAuthenticate(authenticateMiddleware echo.MiddlewareFunc) echo.HandlerFunc {
	return authenticateMiddleware(func(ectx echo.Context) error {
		// ctx := ectx.Request().Context()
		return ectx.JSON(http.StatusOK, ectx.Get("user"))
	})
}

func Login(
	login authenticationservice.LoginFunc,
	createJWT jsonwebtoken.CreateJWTFunc,
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
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := ectx.Bind(&post)
		log.Println("Aa", post)
		if err != nil {
			return ectx.JSON(http.StatusUnauthorized, samasemodels.RESTResponse{Message: "login failed"})
		}
		us, err := login(ctx, post.Email, post.Password)
		if err != nil || us == nil {
			log.Println("Eeh", us, err)
			return ectx.JSON(http.StatusUnauthorized, samasemodels.RESTResponse{Message: "login failed"})
		}
		// log.Println("LOL", us.Email)
		// if us.Email.Verified == false {
		// 	return ectx.JSON(http.StatusOK, struct {
		// 		Verified bool `json:"verified"`
		// 	}{Verified: false})
		// }

		sajwtcl := jsonwebtoken.SamaseJWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenDuration).Unix(),
			},
			User: *us,
		}
		tokenStr, err := createJWT(sajwtcl)
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

func getTokenSource() {
	// config := oauth2.Config{
	// 	ClientID:     "744967159273-rhjtp67vu4075un8hftrp5silgbh2n6f.apps.googleusercontent.com",
	// 	ClientSecret: "Q3LFW6jWMrCyaSB8r7h20fII",
	// 	Endpoint:     google.Endpoint,
	// }
	// config.
}

func verifyIdToken(ctx context.Context, idToken string) {
	// oauth2Service, err := oauth2.NewService(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// tokenInfoCall := oauth2Service.Tokeninfo()
	// tokenInfoCall.IdToken(idToken)

	// tokenInfo, err := tokenInfoCall.Do()
	// if err != nil {
	// 	log.Println(err)
	// }
	// 	return nil, err
	// log.Println(oauth2Service.Userinfo)

	// return tokenInfo, nil
}

func GoogleLogin(
	createUser userservice.CreateUserFunc,
	getUserByEmail userservice.GetUserByEmailFunc,
	verifyIDToken authenticationservice.GoogleVerifyIDTokenFunc,
	createJWT jsonwebtoken.CreateJWTFunc,
	tokenDuration time.Duration,
	sendEmail samasemailservice.SendEmailFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		var post struct {
			IDToken     string `json:"id_token"`
			AccessToken string `json:"access_token"`
		}
		err := ectx.Bind(&post)
		if err != nil {
			return ectx.JSON(http.StatusUnauthorized, nil)
		}
		payload, err := verifyIDToken(ctx, post.IDToken)
		if err != nil {
			return ectx.JSON(http.StatusUnauthorized, nil)
		}
		email := fmt.Sprint(payload.Claims["email"])

		us, err := getUserByEmail(ctx, email)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}

		if us == nil {
			fullname := fmt.Sprint(payload.Claims["name"])
			name := strings.ToLower(strings.Replace(fullname, " ", "", -1))
			log.Println(payload.Claims)
			us = &user.User{
				Name:     name,
				Fullname: fullname,
				Email: &useremail.UserEmail{
					Value:    email,
					Verified: true,
				},
			}
			*us, err = createUser(ctx, *us)
			if err != nil {
				return ectx.JSON(http.StatusInternalServerError, nil)
			}

			// err = sendEmail(ctx, []string{email}, "Ahlan wa sahlan.", "<h1>Ahlan wa sahlan</h1>")
			// if err != nil {
			// 	return ectx.JSON(http.StatusInternalServerError, nil)
			// }
		}

		sajwtcl := jsonwebtoken.SamaseJWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenDuration).Unix(),
			},
			User: *us,
		}
		tokenStr, err := createJWT(sajwtcl)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}
		return ectx.JSON(http.StatusOK, tokenStr)
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

func RefreshToken(
	refreshToken authenticationservice.RefreshTokenFunc,
	tokenDuration time.Duration,
	blackListJWT jsonwebtokenrepo.BlackListJWTFunc,
) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		ctx := ectx.Request().Context()
		authorizationBearer := ectx.Request().Header.Get("Authorization")
		if !strings.Contains(authorizationBearer, "Bearer ") {
			return ectx.JSON(http.StatusUnauthorized, nil)
		}
		token := strings.Replace(authorizationBearer, "Bearer ", "", 1)

		newToken, err := refreshToken(ctx, token, tokenDuration)
		if err != nil {
			log.Println(err)
			return ectx.JSON(http.StatusInternalServerError, nil)
		}

		err = blackListJWT(token)
		if err != nil {
			return ectx.JSON(http.StatusInternalServerError, nil)
		}

		return ectx.JSON(http.StatusOK, struct {
			Token string `json:"token"`
		}{Token: newToken})
	}
}
