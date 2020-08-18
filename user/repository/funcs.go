package userrepo

import (
	"context"
	"fifentory/options"
	"samase/user"
)

type CreateUserFunc func(ctx context.Context, u user.User) (user.User, error)
type GetUsersFunc func(ctx context.Context, opts *options.Options) ([]user.User, error)

type UserFetcher interface {
	GetUsers(ctx context.Context, opts *options.Options) ([]user.User, error)
	WithEmail()
	WithPassword()
}
