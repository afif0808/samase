package userrepo

import (
	"context"
	"fifentory/options"
	"samase/user"
	"time"
)

type CreateUserFunc func(ctx context.Context, u user.User) (user.User, error)
type GetUsersFunc func(ctx context.Context, opts *options.Options) ([]user.User, error)
type UpdateUsersFunc func(ctx context.Context, us user.User, fts []options.Filter) error
type DeleteUsersFunc func(ctx context.Context, fts []options.Filter) error
type SaveEmailConfirmationCodeFunc func(ctx context.Context, code string, expireTime time.Duration) error

type UserFetcher interface {
	GetUsers(ctx context.Context, opts *options.Options) ([]user.User, error)
	WithEmail()
	WithPassword()
	WithPoint()
}
