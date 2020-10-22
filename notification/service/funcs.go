package notificationservice

import (
	"context"
	"samase/notification"
)

type GetNotificationsByUserIDFunc func(ctx context.Context, userID int64) ([]notification.Notification, error)
type CreateNotificationForAllUsersFunc func(ctx context.Context, title string, message string) error
