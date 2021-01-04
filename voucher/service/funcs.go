package voucherservice

import (
	"context"
	"mime/multipart"
	"samase/voucher"
)

type GetVouchersFunc func(ctx context.Context) ([]voucher.Voucher, error)

type SaveVoucherImageFunc func(dest string, f multipart.File) error

type CreateVoucherFunc func(ctx context.Context, vo voucher.Voucher, img multipart.File) (voucher.Voucher, error)
