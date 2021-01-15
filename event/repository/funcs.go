package eventrepo

import (
	"context"
	"fifentory/options"
	"samase/event"
)

type GetEventsFunc func(ctx context.Context, opts *options.Options) ([]event.Event, error)
type CreateEventFunc func(ctx context.Context, ev event.Event) (event.Event, error)
type UpdateEventsFunc func(ctx context.Context, ev event.Event, fts []options.Filter) error
type DeleteEventsFunc func(ctx context.Context, fts []options.Filter) error
