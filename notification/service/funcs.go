package notificationservice

import (
	"context"
	"samase/notification"

	"firebase.google.com/go/messaging"
)

type GetNotificationsByUserIDFunc func(ctx context.Context, userID int64) ([]notification.Notification, error)
type CreateNotificationForAllUsersFunc func(ctx context.Context, notf notification.Notification) error
type MarkNotificationAsReadByIDFunc func(ctx context.Context, id int64) error
type GetUnreadNotificationCountByUserIDFunc func(ctx context.Context, userID int64) (int, error)
type SendWelcomeNotificationFunc func(ctx context.Context, userID int64) error
type SendFirebaseNotificationFunc func(ctx context.Context, msg messaging.Message) error
