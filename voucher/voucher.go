package voucher

type Voucher struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
}
