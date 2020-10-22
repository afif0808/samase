package notificationsqlrepo

import (
	"context"
	"database/sql"
	"log"
	"samase/notification"
	notificationrepo "samase/notification/repository"
)

const (
	notificationTable       = "notification"
	notificationFields      = notificationTable + ".id," + notificationTable + ".name," + notificationTable + ".message," + notificationTable + ".date"
	createNotificationQuery = "INSERT " + notificationTable + " SET notification.name = ? , notification.message = ? , notification.user_id = ?"
)

func CreateNotification(conn *sql.DB) notificationrepo.CreateNotificationFunc {
	return func(ctx context.Context, notf notification.Notification) (notification.Notification, error) {
		res, err := conn.ExecContext(ctx, createNotificationQuery, notf.Name, notf.Message, notf.UserID)
		if err != nil {
			log.Println(err)
			return notf, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			log.Println(err)
			return notf, err
		}
		notf.ID = id
		return notf, nil
	}
}
