package userservice

import (
	"context"
	"samase/user"
)

type DoesNameExistFunc func(ctx context.Context, name string) (bool, error)

type GetUserByEmailFunc func(ctx context.Context, email string) (*user.User, error)
type GetUserByIDFunc func(ctx context.Context, id int64) (*user.User, error)

type GetAllUsersFunc func(ctx context.Context) ([]user.User, error)

type CreateUserFunc func(ctx context.Context, us user.User) (user.User, error)

type UpdateUserFunc func(ctx context.Context, us user.User) error

type SendUserConfirmationEmailFunc func(ctx context.Context, email string) error

type ConfirmUserEmailFunc func(ctx context.Context, email string, code string) error
