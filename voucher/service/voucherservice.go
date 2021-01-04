package voucherservice

import (
	"context"
	"mime/multipart"
	"samase/voucher"
	voucherrepo "samase/voucher/repository"
)

func GetVouchers(
	getVouchers voucherrepo.GetVouchersFunc,
) GetVouchersFunc {
	return func(ctx context.Context) ([]voucher.Voucher, error) {
		vos, err := getVouchers(ctx, nil)
		if err != nil {
			return nil, err
		}
		return vos, nil
	}
}

func CreateVoucher(
	createVoucher voucherrepo.CreateVoucherFunc,
	saveImg SaveVoucherImageFunc,
	imgDest, imgLink string,
) CreateVoucherFunc {
	return func(ctx context.Context, vo voucher.Voucher, img multipart.File) (voucher.Voucher, error) {
		err := saveImg(imgDest, img)
		if err != nil {
			return vo, err
		}
		vo.Image = imgLink + ""
		vo, err = createVoucher(ctx, vo)
		if err != nil {
			return vo, err
		}
		return vo, nil
	}
}
