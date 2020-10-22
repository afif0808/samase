package user

import (
	"samase/useremail"
	"samase/userpassword"
	"samase/userpoint"
	"samase/voucher"
)

type User struct {
	ID       int64                      `json:"id"`
	Name     string                     `json:"name"`
	Fullname string                     `json:"fullname"`
	Email    *useremail.UserEmail       `json:"email,omitempty"`
	Password *userpassword.UserPassword `json:"-"`
	Point    *userpoint.UserPoint       `json:"point,omitempty"`
	Vouchers *voucher.Voucher           `json:"vouchers,omitempty"`
}

type GoogleUser struct {
	Name  string
	Email string
}
