package userpasswordservice

import (
	"context"
	"samase/userpassword"
)

type UpdateUserPasswordFunc func(ctx context.Context, uspa userpassword.UserPassword) error

type GetUserPasswordByUserIDFunc func(ctx context.Context, id int64) (*userpassword.UserPassword, error)
