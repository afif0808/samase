package vouchersqlrepo

import (
	"context"
	"database/sql"
	"log"
	"samase/voucher"
	voucherrepo "samase/voucher/repository"
)

const (
	voucherTable       = "voucher"
	createVoucherQuery = "INSERT " + voucherTable + " SET  voucher.name = ? "
)

func CreateVoucher(conn *sql.DB) voucherrepo.CreateVoucherFunc {
	return func(ctx context.Context, vo voucher.Voucher) (voucher.Voucher, error) {
		res, err := conn.ExecContext(ctx, createVoucherQuery, vo.Name)
		if err != nil {
			log.Println(err)
			return vo, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
		}
		vo.ID = id
		return vo, nil
	}
}
