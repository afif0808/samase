package vouchersqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/voucher"
	voucherrepo "samase/voucher/repository"

	"gorm.io/gorm"
)

const (
	voucherTable        = "voucher"
	fields              = "id,name,image,description"
	createVoucherQuery  = "INSERT " + voucherTable + " SET  voucher.name = ? , voucher.image = ? , voucher.description = ? "
	getVouchersQuery    = "SELECT " + fields + " FROM " + voucherTable
	deleteVouchersQuery = "DELETE FROM " + voucherTable + " "
)

func CreateVoucher(conn *sql.DB) voucherrepo.CreateVoucherFunc {
	return func(ctx context.Context, vo voucher.Voucher) (voucher.Voucher, error) {
		res, err := conn.ExecContext(ctx, createVoucherQuery, vo.Name, vo.Image, vo.Description)
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

func DeleteVouchers(conn *sql.DB) voucherrepo.DeleteVouchersFunc {
	return func(ctx context.Context, fts []options.Filter) error {
		filtersQuery, filtersArgs := options.ParseFiltersToSQLQuery(fts)
		query := deleteVouchersQuery + " " + filtersQuery
		_, err := conn.QueryContext(ctx, query, filtersArgs...)
		if err != nil {
			log.Println(err)
		}
		return nil
	}
}

func UpdateVouchers(db *gorm.DB) voucherrepo.UpdateVouchersFunc {
	return func(ctx context.Context, vo voucher.Voucher, fts []options.Filter) error {
		filtersQuery, filtersArgs := options.GORMParseFiltersToSQLQuery(fts)
		tx := db.WithContext(ctx).
			Table("voucher").
			Where(filtersQuery, filtersArgs[0]).
			Updates(vo)
		err := tx.Error
		if err != nil {
			log.Println(err)
		}
		return err
	}
}
