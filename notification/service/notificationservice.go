package notificationservice

import (
	"context"
	"fifentory/options"
	"samase/notification"
	notificationrepo "samase/notification/repository"
	userrepo "samase/user/repository"
)

func GetNotificationsByUserID(nf notificationrepo.NotificationFetcher) GetNotificationsByUserIDFunc {
	return func(ctx context.Context, userID int64) ([]notification.Notification, error) {
		opts := options.Options{
			Filters: []options.Filter{
				options.Filter{
					By:       "notification.user_id",
					Operator: "=",
					Value:    userID,
				},
			},
		}

		return nf.Fetch(ctx, &opts)
	}
}

func CreateNotificationForAllUsers(
	usf userrepo.UserFetcher,
	createNotification notificationrepo.CreateNotificationFunc,
) CreateNotificationForAllUsersFunc {
	return func(ctx context.Context, title, message string) error {
		uss, err := usf.GetUsers(ctx, nil)
		if err != nil {
			return err
		}
		notf := notification.Notification{
			Name:    title,
			Message: message,
		}
		for _, us := range uss {
			notf.UserID = us.ID
			_, err := createNotification(ctx, notf)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
