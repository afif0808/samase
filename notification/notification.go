package notification

import "time"

type Notification struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
	UserID  int64     `json:"-"`
	Image   string    `json:"image"`
	IsRead  bool      `json:"is_read"`
}
