package notificationrepo

import (
	"context"
	"fifentory/options"
	"samase/notification"
)

type CreateNotificationFunc func(ctx context.Context, notf notification.Notification) (notification.Notification, error)
type UpdateNotificationFunc func(ctx context.Context, notf notification.Notification) (notification.Notification, error)
type NotificationFetcher interface {
	Fetch(ctx context.Context, opts options.Options)
}
