package authenticationmiddleware

import (
	"log"
	"net/http"
	authenticationservice "samase/authentication/service"
	"samase/jsonwebtoken"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
)

func Authenticate(
	parseJWT jsonwebtoken.ParseJWTFunc,
	isLoggedOut authenticationservice.IsLoggedOutFunc,
	secretKey interface{},
	jwtsm jwt.SigningMethod,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			authorizationBearer := ectx.Request().Header.Get("Authorization")
			if !strings.Contains(authorizationBearer, "Bearer ") {
				return ectx.JSON(http.StatusUnauthorized, nil)
			}
			token := strings.Replace(authorizationBearer, "Bearer ", "", 1)
			cl, err := parseJWT(token, secretKey, jwtsm)
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
			return next(ectx)
		}
	}
}
