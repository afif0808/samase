package authenticationmiddleware

import (
	"log"
	"net/http"
	"os"
	authenticationservice "samase/authentication/service"
	"samase/jsonwebtoken"
	jsonwebtokenredisrepo "samase/jsonwebtoken/repository/redis"
	jsonwebtokenservice "samase/jsonwebtoken/service"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"

	"github.com/labstack/echo"
)

func Authenticate(
	parseJWT jsonwebtoken.ParseJWTFunc,
	isLoggedOut authenticationservice.IsLoggedOutFunc,
	jwtconf jsonwebtoken.JWTConfig,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			authorizationBearer := ectx.Request().Header.Get("Authorization")
			if !strings.Contains(authorizationBearer, "Bearer ") {
				return ectx.JSON(http.StatusUnauthorized, nil)
			}
			token := strings.Replace(authorizationBearer, "Bearer ", "", 1)
			cl, err := parseJWT(token, jwtconf.SecretKey, jwtconf.SigningMethod)
			if err != nil {
				log.Println(err)
				return ectx.JSON(http.StatusUnauthorized, nil)
			}
			err = cl.Valid()
			if err != nil {
				return ectx.JSON(http.StatusUnauthorized, nil)
			}

			loggedOut, err := isLoggedOut(token)
			if err != nil || loggedOut {
				return ectx.JSON(http.StatusUnauthorized, nil)
			}
			mcl := cl.(jwt.MapClaims)
			ectx.Set("user", mcl["user"])
			return next(ectx)
		}
	}
}
func redisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   50,
		MaxActive: 10000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", ":6379")
			// Connection error handling
			if err != nil {
				log.Printf("ERROR: fail initializing the redis pool: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}
func InjectAuthenticate() echo.MiddlewareFunc {
	rp := redisPool()
	rc, _ := rp.Dial()

	parseJWT := jsonwebtokenservice.ParseJWT()
	isJWTBlackListed := jsonwebtokenredisrepo.IsJWTBlackListed(rc)
	isLoggedOut := authenticationservice.IsLoggedOut(isJWTBlackListed)
	jwtconf := jsonwebtoken.JWTConfig{
		SecretKey:     []byte("itssignaturekey"),
		SigningMethod: jwt.SigningMethodHS256,
	}
	return Authenticate(parseJWT, isLoggedOut, jwtconf)
}
