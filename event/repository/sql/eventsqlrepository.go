package eventsqlrepository

import (
	"context"
	"database/sql"
	"fifentory/options"
	"log"
	"samase/event"
	eventrepo "samase/event/repository"

	"gorm.io/gorm"
)

const (
	fields         = "id,name,image,description,started_at,ended_at,created_at"
	tableName      = "event"
	getEventsQuery = "SELECT " + fields + " FROM " + tableName + " "
)

func GetEvents(conn *sql.DB) eventrepo.GetEventsFunc {
	return func(ctx context.Context, opts *options.Options) ([]event.Event, error) {
		optionsQuery, optionsArgs := options.ParseOptionsToSQLQuery(opts)
		query := getEventsQuery + " " + optionsQuery
		// log.Println(query)
		rows, err := conn.QueryContext(ctx, query, optionsArgs...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer rows.Close()
		evs := []event.Event{}
		for rows.Next() {
			ev := event.Event{}
			err = rows.Scan(&ev.ID, &ev.Name, &ev.Image, &ev.Description, &ev.StartedAt, &ev.EndedAt, &ev.CreatedAt)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			evs = append(evs, ev)
		}
		return evs, nil
	}
}
func CreateEvent(db *gorm.DB) eventrepo.CreateEventFunc {
	return func(ctx context.Context, ev event.Event) (event.Event, error) {
		res := db.Table("event").Create(&ev)
		err := res.Error
		if err != nil {
			log.Println(err)
		}
		return ev, err
	}
}

func UpdateEvents(db *gorm.DB) eventrepo.UpdateEventsFunc {
	return func(ctx context.Context, ev event.Event, fts []options.Filter) error {
		filtersQuery, filtersArgs := options.GORMParseFiltersToSQLQuery(fts)
		res := db.Table("event").
			Where(filtersQuery, filtersArgs...).
			Updates(ev)
		err := res.Error
		if err != nil {
			log.Println(err)
		}
		return err
	}
}

func DeleteEvents(db *gorm.DB) eventrepo.DeleteEventsFunc {
	return func(ctx context.Context, fts []options.Filter) error {
		filtersQuery, filtersArgs := options.GORMParseFiltersToSQLQuery(fts)
		res := db.Table("event").Where(filtersQuery, filtersArgs...).Delete(&event.Event{})
		err := res.Error
		if err != nil {
			log.Println(err)
		}
		return err
	}
}
