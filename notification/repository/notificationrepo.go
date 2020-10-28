package notificationrepo

import (
	"context"
	"fifentory/options"
	"samase/notification"
)

type GetNotificationsFunc func(ctx context.Context, opts *options.Options) ([]notification.Notification, error)
type CreateNotificationFunc func(ctx context.Context, notf notification.Notification) (notification.Notification, error)
type UpdateNotificationFunc func(ctx context.Context, notf notification.Notification) error
type NotificationFetcher interface {
	Fetch(ctx context.Context, opts *options.Options) ([]notification.Notification, error)
}
type MarkNotificationAsReadFunc func(ctx context.Context, fts []options.Filter) error

type GetNotificationFetcherFunc func() NotificationFetcher
