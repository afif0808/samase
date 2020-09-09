package notification

import "time"

type Notification struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Message string    `json:"message"`
	Date    time.Time `json:"date"`
}
