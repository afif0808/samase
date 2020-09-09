package uservoucherjunctionrepo

import (
	"context"
	"fifentory/options"
	"samase/uservoucherjunction"
)

type CreateUserVoucherJunctionFunc func(ctx context.Context, usvo uservoucherjunction.UserVoucherJunction) (uservoucherjunction.UserVoucherJunction, error)

type UserFetcher interface {
	GetUserVoucherJunctions(ctx context.Context, opts *options.Options) ([]uservoucherjunction.UserVoucherJunction, error)
}
