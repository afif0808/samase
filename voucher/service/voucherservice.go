package voucherservice

import (
	"context"
	"fifentory/options"
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

func DeleteVoucherByID(deleteVoucher voucherrepo.DeleteVouchersFunc) DeleteVoucherByIDFunc {
	return func(ctx context.Context, id int64) error {
		fts := []options.Filter{
			options.Filter{
				By:       "id",
				Operator: "=",
				Value:    id,
			},
		}
		err := deleteVoucher(ctx, fts)
		return err
	}
}
