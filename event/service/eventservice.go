package eventservice

import (
	"context"
	"fifentory/options"
	"samase/event"
	eventrepo "samase/event/repository"
	"strings"
	"time"
)

func GetEvents(getEvents eventrepo.GetEventsFunc) GetEventsFunc {
	return func(ctx context.Context, keyword string) ([]event.Event, error) {
		evs, err := getEvents(ctx, nil)
		if err != nil {
			return nil, err
		}
		return evs, err
	}
}

func GetOngoingEvents(getEvents eventrepo.GetEventsFunc) GetOngoingEventsFunc {
	return func(ctx context.Context, keyword string) ([]event.Event, error) {
		today := time.Now().String()
		today = strings.Split(today, " ")[0]
		opts := options.Options{Filters: []options.Filter{
			options.Filter{
				Operator: "<=",
				Value:    today,
				By:       "started_at",
			},
			options.Filter{
				Operator: ">=",
				Value:    today,
				By:       "ended_at",
			},
		}}
		evs, err := getEvents(ctx, &opts)
		if err != nil {
			return nil, err
		}
		return evs, nil
	}
}

func CreateEvent(createEvent eventrepo.CreateEventFunc) CreateEventFunc {
	return func(ctx context.Context, ev event.Event) (event.Event, error) {
		return createEvent(ctx, ev)
	}
}
func UpdateEventByID(updateEvents eventrepo.UpdateEventsFunc) UpdateEventByIDFunc {
	return func(ctx context.Context, ev event.Event) error {
		fts := []options.Filter{
			options.Filter{
				Operator: "=",
				By:       "id",
				Value:    ev.ID,
			},
		}
		return updateEvents(ctx, ev, fts)
	}
}

func DeleteEventByID(deleteEvents eventrepo.DeleteEventsFunc) DeleteEventByIDFunc {
	return func(ctx context.Context, id int64) error {
		fts := []options.Filter{
			options.Filter{
				Operator: "=",
				By:       "id",
				Value:    id,
			},
		}
		return deleteEvents(ctx, fts)
	}
}
