package uservoucherjunction

type UserVoucherJunction struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	VoucherID int64 `json:"voucher_id"`
}
