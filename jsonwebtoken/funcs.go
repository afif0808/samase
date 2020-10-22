package jsonwebtoken

import (
	"github.com/dgrijalva/jwt-go"
)

type CreateJWTFunc func(jwtcl jwt.Claims) (string, error)
type ParseJWTFunc func(tokenStr string) (jwt.Claims, error)

// type ParseJWTFunc func(tokenStr string, secrectKey interface{}, jwtsm jwt.SigningMethod) (jwt.Claims, error)
