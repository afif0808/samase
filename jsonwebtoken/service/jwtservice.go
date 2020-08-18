package jsonwebtokenservice

import (
	"errors"
	"log"
	"samase/jsonwebtoken"

	"github.com/dgrijalva/jwt-go"
)

func CreateJWT() jsonwebtoken.CreateJWTFunc {
	return func(jwtcl jwt.Claims, secretkey interface{}, jwtsm jwt.SigningMethod) (string, error) {
		token := jwt.NewWithClaims(jwtsm, jwtcl)
		tokenStr, err := token.SignedString(secretkey)
		if err != nil {
			log.Println(err)
			return "", err
		}
		return tokenStr, nil
	}
}

func ParseJWT() jsonwebtoken.ParseJWTFunc {
	return func(tokenStr string, secrectKey interface{}, jwtsm jwt.SigningMethod) (jwt.Claims, error) {
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Signing method invalid")
			} else if method != jwtsm {
				return nil, errors.New("Signing method invalid")
			}
			return secrectKey, nil
		})
		if err != nil {
			log.Println()
			return nil, err
		}
		return token.Claims, nil
	}
}
