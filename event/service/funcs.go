package eventservice

import (
	"context"
	"samase/event"
)

type GetEventsFunc func(ctx context.Context, keyword string) ([]event.Event, error)
type GetOngoingEventsFunc func(ctx context.Context, keyword string) ([]event.Event, error)
type CreateEventFunc func(ctx context.Context, ev event.Event) (event.Event, error)
type UpdateEventByIDFunc func(ctx context.Context, ev event.Event) error
type DeleteEventByIDFunc func(ctx context.Context, id int64) error
