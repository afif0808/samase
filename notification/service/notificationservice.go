package notificationservice

import (
	"context"
	"fifentory/options"
	"log"
	"samase/notification"
	notificationrepo "samase/notification/repository"
	userrepo "samase/user/repository"
	userservice "samase/user/service"
	"time"

	"firebase.google.com/go/messaging"
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
	sendFirebaseNotification SendFirebaseNotificationFunc,
	getUserWSs userservice.GetUserWSsFunc,
) CreateNotificationForAllUsersFunc {
	return func(ctx context.Context, notf notification.Notification) error {
		uss, err := usf.GetUsers(ctx, nil)
		if err != nil {
			return err
		}
		for _, us := range uss {
			notf.UserID = us.ID
			_, err := createNotification(ctx, notf)
			if err != nil {
				return err
			}
		}
		msg := messaging.Message{
			Notification: &messaging.Notification{
				Title:    notf.Name,
				Body:     notf.Message,
				ImageURL: notf.Image,
			},
			Topic: "topic",
		}
		err = sendFirebaseNotification(ctx, msg)
		if err != nil {
			return err
		}

		userWSs := getUserWSs()
		for ws := range userWSs {
			err = ws.WriteJSON(notf)
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

func SendWelcomeNotification(createNotification notificationrepo.CreateNotificationFunc) SendWelcomeNotificationFunc {
	return func(ctx context.Context, userID int64) error {
		notf := notification.Notification{
			UserID:  userID,
			Date:    time.Now(),
			Name:    "Selamat , anda telah menajdi member Samase",
			Message: "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen bookIt has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
		}
		notf, err := createNotification(ctx, notf)
		return err
	}
}

func SendFirebaseNotification(msgcl *messaging.Client) SendFirebaseNotificationFunc {
	return func(ctx context.Context, msg messaging.Message) error {
		_, err := msgcl.Send(ctx, &msg)
		if err != nil {
			log.Println(err)
		}
		return err
	}
}
