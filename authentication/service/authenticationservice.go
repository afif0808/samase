package authenticationservice

import (
	"context"
	"errors"
	"fifentory/options"
	"fmt"
	"samase/jsonwebtoken"
	jsonwebtokenrepo "samase/jsonwebtoken/repository"
	"samase/user"
	userrepo "samase/user/repository"
	userservice "samase/user/service"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
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
func GoogleVerifyIDToken(audience string, credentialsFile string) GoogleVerifyIDTokenFunc {
	return func(ctx context.Context, IDToken string) (*idtoken.Payload, error) {
		validator, err := idtoken.NewValidator(ctx, idtoken.WithCredentialsFile(credentialsFile))
		if err != nil {
			return nil, err
		}
		pl, err := validator.Validate(ctx, IDToken, audience)

		if err != nil {
			return nil, err
		}
		return pl, nil
	}
}

func checkPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(usfe userrepo.UserFetcher) LoginFunc {
	return func(ctx context.Context, email, password string) (*user.User, error) {
		usfe.WithEmail()
		usfe.WithPassword()
		uss, err := usfe.GetUsers(ctx, &options.Options{
			Filters: []options.Filter{
				options.Filter{
					By:       "user_email.value",
					Operator: "=",
					Value:    email,
				},
			},
		})
		if err != nil {
			return nil, err
		}
		if len(uss) <= 0 {
			return nil, errors.New("No user with that name exists")
		}

		us := uss[0]
		if !checkPasswordHash(us.Password.Hash, password) {
			return nil, err
		}
		return &us, nil
	}
}

func RefreshToken(
	createJWT jsonwebtoken.CreateJWTFunc,
	parseJWT jsonwebtoken.ParseJWTFunc,
	getUserByID userservice.GetUserByIDFunc,
) RefreshTokenFunc {
	return func(ctx context.Context, token string, tokenDuration time.Duration) (string, error) {
		cl, err := parseJWT(token)
		if err != nil {
			return "", err
		}
		if err := cl.Valid(); err != nil {
			return "", err
		}

		mapcl := cl.(jwt.MapClaims)
		usmap := mapcl["user"].(map[string]interface{})
		userID, err := strconv.ParseInt(fmt.Sprint(usmap["id"]), 10, 64)
		if err != nil {
			return "", err
		}
		us, err := getUserByID(ctx, userID)
		if err != nil {
			return "", err
		}
		sajwtcl := jsonwebtoken.SamaseJWTClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenDuration).Unix(),
			},
			User: *us,
		}
		newToken, err := createJWT(sajwtcl)
		if err != nil {
			return "", err
		}
		return newToken, nil
	}
}
