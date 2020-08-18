package userpassword

type UserPassword struct {
	UserID int64  `json:"-"`
	Hash   string `json:"-"`
	Value  string `json:"value"`
}
