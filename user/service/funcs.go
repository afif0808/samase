package userservice

import (
	"context"
	"samase/user"

	"github.com/gorilla/websocket"
)

type DoesNameExistFunc func(ctx context.Context, name string) (bool, error)

type GetUserByEmailFunc func(ctx context.Context, email string) (*user.User, error)

type GetUserByIDFunc func(ctx context.Context, id int64) (*user.User, error)

type GetAllUsersFunc func(ctx context.Context) ([]user.User, error)

type CreateUserFunc func(ctx context.Context, us user.User) (user.User, error)

type UpdateUserFunc func(ctx context.Context, us user.User) error

type SendUserConfirmationEmailFunc func(ctx context.Context, email string) error

type ConfirmUserEmailFunc func(ctx context.Context, email string, code string) error

type SendPasswordRecoveryCodeFunc func(ctx context.Context, email string) error

type ConfirmPasswordRecoveryCodeFunc func(ctx context.Context, email, code string) error

type SendAccountPasswordRecoveryLinkFunc func(ctx context.Context, email string) error

type GetUserIDByCodeFunc func(ctx context.Context, code string) (int64, error)

type RecoverUserPasswordFunc func(ctx context.Context, code, password string) error

type RegisterUserWebSocketFunc func(ws *websocket.Conn) error
type UnregisterUserWebSocketFunc func(ws *websocket.Conn) error
type GetUserWSsFunc func() map[*websocket.Conn]struct{}