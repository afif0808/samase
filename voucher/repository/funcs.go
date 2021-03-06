package voucherrepo

import (
	"context"
	"fifentory/options"
	"samase/voucher"
)

type CreateVoucherFunc func(ctx context.Context, vo voucher.Voucher) (voucher.Voucher, error)

type GetVouchersFunc func(ctx context.Context, opts *options.Options) ([]voucher.Voucher, error)

type DeleteVouchersFunc func(ctx context.Context, fts []options.Filter) error

type UpdateVouchersFunc func(ctx context.Context, vo voucher.Voucher, fts []options.Filter) error

type VoucherFetcher interface {
	GetVouchers(ctx context.Context, opts *options.Options) ([]voucher.Voucher, error)
}
