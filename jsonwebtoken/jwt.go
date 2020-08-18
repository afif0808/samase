package jsonwebtoken

import (
	"samase/user"

	"github.com/dgrijalva/jwt-go"
)

type SamaseJWTClaims struct {
	jwt.StandardClaims
	User user.User
}
