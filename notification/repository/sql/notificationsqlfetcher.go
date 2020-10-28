package notificationsqlrepo

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/notification"
	notificationrepo "samase/notification/repository"
)

type receiver struct {
	Notification *notification.Notification
}

type NotificationSQLFetcher struct {
	joins    string
	scanDest []interface{}
	fields   string
	Receiver *receiver
	conn     *sql.DB
}

func GetNotificationSQLFetcher(conn *sql.DB) notificationrepo.GetNotificationFetcherFunc {
	return func() notificationrepo.NotificationFetcher {
		notfsf := NewNotificationSQLFetcher(conn)
		return &notfsf
	}
}

func NewNotificationSQLFetcher(conn *sql.DB) NotificationSQLFetcher {
	notfsf := NotificationSQLFetcher{
		Receiver: &receiver{Notification: &notification.Notification{}},
		conn:     conn,
	}
	return notfsf
}

func (notfsf *NotificationSQLFetcher) Fetch(ctx context.Context, opts *options.Options) ([]notification.Notification, error) {
	notfsf.fields += notificationFields

	notfsf.scanDest = append(
		notfsf.scanDest,
		&notfsf.Receiver.Notification.ID,
		&notfsf.Receiver.Notification.Name,
		&notfsf.Receiver.Notification.Message,
		&notfsf.Receiver.Notification.Date,
		&notfsf.Receiver.Notification.IsRead,
	)

	defer func() {
		notfsf.fields = ""
		notfsf.joins = ""
		notfsf.scanDest = []interface{}{}
	}()

	optionsQuery, optionsArgs := options.ParseOptionsToSQLQuery(opts)

	query := "SELECT " + notfsf.fields + " FROM " + notificationTable + " " + optionsQuery
	rows, err := notfsf.conn.QueryContext(ctx, query, optionsArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notfs := []notification.Notification{}

	for rows.Next() {
		err := rows.Scan(notfsf.scanDest...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if notfsf.Receiver.Notification != nil {
			notf := notification.Notification{
				ID:      notfsf.Receiver.Notification.ID,
				Name:    notfsf.Receiver.Notification.Name,
				Message: notfsf.Receiver.Notification.Message,
				Date:    notfsf.Receiver.Notification.Date,
				IsRead:  notfsf.Receiver.Notification.IsRead,
			}
			notfs = append(notfs, notf)
		}
	}

	return notfs, nil
}
