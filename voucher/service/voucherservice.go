package voucherservice

import (
	"context"
	"fifentory/options"
	"samase/voucher"
	voucherrepo "samase/voucher/repository"
)

func GetVouchers(
	getVouchers voucherrepo.GetVouchersFunc,
) GetVouchersFunc {
	return func(ctx context.Context, keyword string) ([]voucher.Voucher, error) {
		opts := options.Options{
			Filters: []options.Filter{
				options.Filter{
					Operator: "LIKE",
					Value:    keyword,
					By:       "name",
				},
			},
		}
		vos, err := getVouchers(ctx, &opts)
		if err != nil {
			return nil, err
		}
		return vos, nil
	}
}

func CreateVoucher(
	createVoucher voucherrepo.CreateVoucherFunc,
) CreateVoucherFunc {
	return func(ctx context.Context, vo voucher.Voucher) (voucher.Voucher, error) {
		vo, err := createVoucher(ctx, vo)
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

func UpdateVoucherByID(updateVouchers voucherrepo.UpdateVouchersFunc) UpdateVoucherByIDFunc {
	return func(ctx context.Context, vo voucher.Voucher) error {
		fts := []options.Filter{options.Filter{
			Operator: "=",
			By:       "id",
			Value:    vo.ID,
		}}

		return updateVouchers(ctx, vo, fts)
	}
}
