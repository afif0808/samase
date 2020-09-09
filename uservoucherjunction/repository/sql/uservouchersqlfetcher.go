package uservoucherjunctionsqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/uservoucherjunction"
)

type receiver struct {
	UserVoucherJunction *uservoucherjunction.UserVoucherJunction
}

type UserVoucherJunctionSQLFetcher struct {
	joins    string
	scanDest []interface{}
	fields   string
	Receiver *receiver
	conn     *sql.DB
}

func NewUserVoucherJunctionSQLFetcher(conn *sql.DB) UserVoucherJunctionSQLFetcher {
	usvosf := UserVoucherJunctionSQLFetcher{
		fields: "user_voucher.id,user_voucher.user_id,user_voucher.voucher_id",
		Receiver: &receiver{
			UserVoucherJunction: &uservoucherjunction.UserVoucherJunction{},
		},
		conn: conn,
	}
	usvosf.scanDest = []interface{}{
		&usvosf.Receiver.UserVoucherJunction.ID, &usvosf.Receiver.UserVoucherJunction.UserID, &usvosf.Receiver.UserVoucherJunction.VoucherID,
	}
	return usvosf
}

func (usvosf *UserVoucherJunctionSQLFetcher) Getuservoucherjunctions(ctx context.Context, opts *options.Options) ([]uservoucherjunction.UserVoucherJunction, error) {
	optionsQuery, optionsArgs := options.ParseOptionsToSQLQuery(opts)
	rows, err := usvosf.conn.QueryContext(ctx, "SELECT "+usvosf.fields+" FROM "+uservoucherjunctionTable+" "+usvosf.joins+" "+optionsQuery, optionsArgs...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	usvos := []uservoucherjunction.UserVoucherJunction{}
	for rows.Next() {
		rows.Scan(usvosf.scanDest...)

		usvo := uservoucherjunction.UserVoucherJunction{
			ID:        usvosf.Receiver.UserVoucherJunction.ID,
			UserID:    usvosf.Receiver.UserVoucherJunction.UserID,
			VoucherID: usvosf.Receiver.UserVoucherJunction.VoucherID,
		}
		usvos = append(usvos, usvo)
	}
	return usvos, nil
}
