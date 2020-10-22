package useremail

type UserEmail struct {
	Value    string `json:"value"`
	UserID   int64  `json:"-"`
	Verified bool   `json:"verified"`
}
