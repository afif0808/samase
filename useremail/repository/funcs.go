package useremailrepo

import (
	"context"
	"fifentory/options"
	"samase/useremail"
)

type CreateUserEmailFunc func(ctx context.Context, usem useremail.UserEmail) (useremail.UserEmail, error)
type GetUserEmailsFunc func(ctx context.Context, opts *options.Options) ([]useremail.UserEmail, error)
type UpdateUserEmailsFunc func(ctx context.Context, usem useremail.UserEmail, fts []options.Filter) error
