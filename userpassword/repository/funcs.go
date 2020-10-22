package userpasswordrepo

import (
	"context"
	"fifentory/options"
	"samase/userpassword"
)

type CreateUserPasswordFunc func(ctx context.Context, uspa userpassword.UserPassword) (userpassword.UserPassword, error)
type UpdateUserPasswordFunc func(ctx context.Context, uspa userpassword.UserPassword, fts []options.Filter) error
type GetUserPasswordsFunc func(ctx context.Context, opts *options.Options) ([]userpassword.UserPassword, error)
