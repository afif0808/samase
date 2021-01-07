package eventservice

import (
	"context"
	"fifentory/options"
	"log"
	"samase/event"
	eventrepo "samase/event/repository"
	"strings"
	"time"
)

func GetEvents(getEvents eventrepo.GetEventsFunc) GetEventsFunc {
	return func(ctx context.Context) ([]event.Event, error) {
		evs, err := getEvents(ctx, nil)
		if err != nil {
			return nil, err
		}
		return evs, err
	}
}

func GetOngoingEvents(getEvents eventrepo.GetEventsFunc) GetOngoingEventsFunc {
	return func(ctx context.Context) ([]event.Event, error) {
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
		log.Println(today)
		evs, err := getEvents(ctx, &opts)
		if err != nil {
			return nil, err
		}
		return evs, nil
	}
}
