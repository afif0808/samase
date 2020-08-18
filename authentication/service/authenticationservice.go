package authenticationservice

import (
	jsonwebtokenrepo "samase/jsonwebtoken/repository"
)

func Logout(
	blackListJWT jsonwebtokenrepo.BlackListJWTFunc,
) LogoutFunc {
	return func(token string) error {
		return blackListJWT(token)
	}
}

func IsLoggedOut(isJWTBlackListed jsonwebtokenrepo.IsJWTBlackListedFunc) IsLoggedOutFunc {
	return func(token string) (bool, error) {
		return isJWTBlackListed(token)
	}
}
