package uservoucherjunctionservice

import (
	"context"
	"samase/uservoucherjunction"
	uservoucherjunctionrepo "samase/uservoucherjunction/repository"
)

func ClaimVoucher(
	createUserVoucherJunction uservoucherjunctionrepo.CreateUserVoucherJunctionFunc,
) ClaimVoucherFunc {
	return func(ctx context.Context, usvo uservoucherjunction.UserVoucherJunction) (uservoucherjunction.UserVoucherJunction, error) {
		usvo, err := createUserVoucherJunction(ctx, usvo)
		return usvo, err
	}
}
