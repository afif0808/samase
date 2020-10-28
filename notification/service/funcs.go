package notificationservice

import (
	"context"
	"samase/notification"
)

type GetNotificationsByUserIDFunc func(ctx context.Context, userID int64) ([]notification.Notification, error)
type CreateNotificationForAllUsersFunc func(ctx context.Context, title string, message string) error
type MarkNotificationAsReadByIDFunc func(ctx context.Context, id int64) error
type GetUnreadNotificationCountByUserIDFunc func(ctx context.Context, userID int64) (int, error)
