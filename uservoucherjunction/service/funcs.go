package uservoucherjunctionservice

import (
	"context"
	"samase/uservoucherjunction"
)

type ClaimVoucherFunc func(ctx context.Context, usvo uservoucherjunction.UserVoucherJunction) (uservoucherjunction.UserVoucherJunction, error)
