package user

import (
	"samase/useremail"
	"samase/userpassword"
)

type User struct {
	ID       int64                      `json:"id"`
	Name     string                     `json:"name"`
	Fullname string                     `json:"fullname"`
	Email    *useremail.UserEmail       `json:"email"`
	Password *userpassword.UserPassword `json:"-"`
}
