package notificationservice

import (
	"context"
	"fifentory/options"
	"samase/notification"
	notificationrepo "samase/notification/repository"
	userrepo "samase/user/repository"
)

func GetNotificationsByUserID(getNotifications notificationrepo.GetNotificationsFunc) GetNotificationsByUserIDFunc {
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

		return getNotifications(ctx, &opts)
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

func MarkNotificationAsReadByID(markNotificationAsRead notificationrepo.MarkNotificationAsReadFunc) MarkNotificationAsReadByIDFunc {
	return func(ctx context.Context, id int64) error {
		fts := []options.Filter{
			options.Filter{
				By:       "notification.id",
				Operator: "=",
				Value:    id,
			},
		}
		return markNotificationAsRead(ctx, fts)
	}
}

func GetUnreadNotificationsByUserID() {

}

func GetUnreadNotificationCountByUserID(getNotifications notificationrepo.GetNotificationsFunc) GetUnreadNotificationCountByUserIDFunc {
	return func(ctx context.Context, userID int64) (int, error) {
		opts := options.Options{
			Filters: []options.Filter{
				options.Filter{
					By:       "notification.user_id",
					Operator: "=",
					Value:    userID,
				},
				options.Filter{
					By:       "notification.is_read",
					Operator: "=",
					Value:    0,
				},
			},
		}
		notfs, err := getNotifications(ctx, &opts)
		if err != nil {
			return 0, err
		}
		return len(notfs), err
	}
}
