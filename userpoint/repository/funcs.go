package userpointrepo

import (
	"context"
	"samase/userpoint"
)

type CreateUserPointFunc func(ctx context.Context, uspo userpoint.UserPoint) (userpoint.UserPoint, error)
