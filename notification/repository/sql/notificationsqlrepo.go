package notificationsqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/notification"
	notificationrepo "samase/notification/repository"
)

const (
	notificationTable           = "notification"
	notificationFields          = notificationTable + ".id," + notificationTable + ".name," + notificationTable + ".message," + notificationTable + ".date , " + notificationTable + ".is_read"
	createNotificationQuery     = "INSERT " + notificationTable + " SET notification.name = ? , notification.message = ? , notification.user_id = ?"
	markNotificationAsReadQuery = "UPDATE " + notificationTable + " SET is_read = 1 "
	getNotificationsQuery       = "SELECT " + notificationFields + " FROM " + notificationTable
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

func GetNotifications(conn *sql.DB) notificationrepo.GetNotificationsFunc {
	return func(ctx context.Context, opts *options.Options) ([]notification.Notification, error) {
		optionsQuery, optionsArgs := options.ParseOptionsToSQLQuery(opts)
		query := getNotificationsQuery + " " + optionsQuery
		rows, err := conn.QueryContext(ctx, query, optionsArgs...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer rows.Close()
		notfs := []notification.Notification{}
		for rows.Next() {
			notf := notification.Notification{}
			err := rows.Scan(&notf.ID, &notf.Name, &notf.Message, &notf.Date, &notf.IsRead)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			notfs = append(notfs, notf)
		}
		return notfs, nil
	}
}

func MarkNotificationAsRead(conn *sql.DB) notificationrepo.MarkNotificationAsReadFunc {
	return func(ctx context.Context, fts []options.Filter) error {
		filtersQuery, filtersArgs := options.ParseFiltersToSQLQuery(fts)
		query := markNotificationAsReadQuery + " " + filtersQuery
		_, err := conn.ExecContext(ctx, query, filtersArgs...)
		if err != nil {
			log.Println(err)
		}
		return err
	}
}
