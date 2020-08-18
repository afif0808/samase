package userpasswordrepo

import (
	"context"
	"samase/userpassword"
)

type CreateUserPasswordFunc func(ctx context.Context, uspa userpassword.UserPassword) (userpassword.UserPassword, error)
