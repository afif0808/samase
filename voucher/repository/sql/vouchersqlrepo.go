package vouchersqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/voucher"
	voucherrepo "samase/voucher/repository"
)

const (
	voucherTable       = "voucher"
	fields             = "id,name,image,description"
	createVoucherQuery = "INSERT " + voucherTable + " SET  voucher.name = ? , voucher.image = ? "
	getVouchersQuery   = "SELECT " + fields + " FROM " + voucherTable
)

func CreateVoucher(conn *sql.DB) voucherrepo.CreateVoucherFunc {
	return func(ctx context.Context, vo voucher.Voucher) (voucher.Voucher, error) {
		res, err := conn.ExecContext(ctx, createVoucherQuery, vo.Name, vo.Image)
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
func GetVouchers(conn *sql.DB) voucherrepo.GetVouchersFunc {
	return func(ctx context.Context, opts *options.Options) ([]voucher.Voucher, error) {
		optionsQuery, optionsArgs := options.ParseOptionsToSQLQuery(opts)
		query := getVouchersQuery + " " + optionsQuery
		rows, err := conn.QueryContext(ctx, query, optionsArgs...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer rows.Close()
		vos := []voucher.Voucher{}
		for rows.Next() {
			vo := voucher.Voucher{}
			err := rows.Scan(&vo.ID, &vo.Name, &vo.Image, &vo.Description)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			vos = append(vos, vo)
		}
		return vos, nil
	}
}
