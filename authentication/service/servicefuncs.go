package authenticationservice

import (
	"context"
	"samase/user"
	"time"

	"google.golang.org/api/idtoken"
)

type LogoutFunc func(token string) error
type IsLoggedOutFunc func(token string) (bool, error)
type GoogleVerifyIDTokenFunc func(ctx context.Context, IDToken string) (*idtoken.Payload, error)
type LoginFunc func(ctx context.Context, email, password string) (*user.User, error)
type RefreshTokenFunc func(ctx context.Context, token string, tokenDuration time.Duration) (string, error)
