package uservoucherjunctionsqlrepo

import (
	"context"
	"database/sql"
	"samase/uservoucherjunction"
	uservoucherjunctionrepo "samase/uservoucherjunction/repository"
)

const (
	uservoucherjunctionTable       = "user_voucher"
	createuservoucherjunctionQuery = "INSERT " + uservoucherjunctionTable + " SET user_voucher.user_id = ? , user_voucher.voucher_id = ?"
)

func Createuservoucherjunction(conn *sql.DB) uservoucherjunctionrepo.CreateUserVoucherJunctionFunc {
	return func(ctx context.Context, usvo uservoucherjunction.UserVoucherJunction) (uservoucherjunction.UserVoucherJunction, error) {
		res, err := conn.ExecContext(ctx, createuservoucherjunctionQuery, usvo.UserID, usvo.VoucherID)
		if err != nil {
			return usvo, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			usvo.ID = id
		}
		return usvo, nil
	}
}
