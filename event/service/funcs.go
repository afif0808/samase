package eventservice

import (
	"context"
	"samase/event"
)

type GetEventsFunc func(ctx context.Context) ([]event.Event, error)
type GetOngoingEventsFunc func(ctx context.Context) ([]event.Event, error)
