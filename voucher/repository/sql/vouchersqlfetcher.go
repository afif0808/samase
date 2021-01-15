package vouchersqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/voucher"
)

type receiver struct {
	Voucher *voucher.Voucher
}

type VoucherSQLFetcher struct {
	joins    string
	scanDest []interface{}
	fields   string
	Receiver *receiver
	conn     *sql.DB
}

func NewVoucherSQLFetcher(conn *sql.DB) VoucherSQLFetcher {
	vosf := VoucherSQLFetcher{
		conn:   conn,
		fields: "voucher.id,voucher.name",
		Receiver: &receiver{
			Voucher: &voucher.Voucher{},
		},
	}
	vosf.scanDest = []interface{}{&vosf.Receiver.Voucher.ID, &vosf.Receiver.Voucher.Name}
	return vosf
}

func (vosf *VoucherSQLFetcher) GetVouchers(ctx context.Context, opts *options.Options) ([]voucher.Voucher, error) {
	optionsQuery, optionsArgs := options.ParseOptionsToSQLQuery(opts)
	rows, err := vosf.conn.QueryContext(ctx, "SELECT "+vosf.fields+" FROM "+voucherTable+" "+vosf.joins+" "+optionsQuery, optionsArgs...)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()
	vos := []voucher.Voucher{}
	for rows.Next() {
		err := rows.Scan(vosf.scanDest...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		vo := voucher.Voucher{
			ID:   vosf.Receiver.Voucher.ID,
			Name: vosf.Receiver.Voucher.Name,
		}
		vos = append(vos, vo)
	}
	return vos, nil
}
