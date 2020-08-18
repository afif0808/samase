package useremail

type UserEmail struct {
	Value    string `json:"value"`
	UserID   int64  `json:"user_id"`
	Verified bool   `json:"verified"`
}
