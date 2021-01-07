package eventrepo

import (
	"context"
	"fifentory/options"
	"samase/event"
)

type GetEventsFunc func(ctx context.Context, opts *options.Options) ([]event.Event, error)
