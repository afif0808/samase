package voucherservice

import (
	"context"
	"mime/multipart"
	"samase/voucher"
)

type GetVouchersFunc func(ctx context.Context, search string) ([]voucher.Voucher, error)

type SaveVoucherImageFunc func(dest string, f multipart.File) error

type CreateVoucherFunc func(ctx context.Context, vo voucher.Voucher) (voucher.Voucher, error)

type DeleteVoucherByIDFunc func(ctx context.Context, id int64) error

type UpdateVoucherByIDFunc func(ctx context.Context, vo voucher.Voucher) error
