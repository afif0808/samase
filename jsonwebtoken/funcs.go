package jsonwebtoken

import (
	"github.com/dgrijalva/jwt-go"
)

type CreateJWTFunc func(jwtcl jwt.Claims, secretkey interface{}, signingMethod jwt.SigningMethod) (string, error)
type ParseJWTFunc func(tokenStr string, secrectKey interface{}, jwtsm jwt.SigningMethod) (jwt.Claims, error)
