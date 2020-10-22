package jsonwebtoken

import (
	"samase/user"

	"github.com/dgrijalva/jwt-go"
)

type SamaseJWTClaims struct {
	jwt.StandardClaims
	User user.User `json:"user"`
}

// JWTConfig consists configuration of jwt
type JWTConfig struct {
	SecretKey     interface{}
	SigningMethod jwt.SigningMethod
}

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}
