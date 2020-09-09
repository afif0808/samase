package voucherrepo

import (
	"context"
	"fifentory/options"
	"samase/voucher"
)

type CreateVoucherFunc func(ctx context.Context, vo voucher.Voucher) (voucher.Voucher, error)

type VoucherFetcher interface {
	GetVouchers(ctx context.Context, opts *options.Options) ([]voucher.Voucher, error)
}
