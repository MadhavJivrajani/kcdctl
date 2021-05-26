package watcher

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

// WatcherLister watches the state of the system
// and lists events generated on change in system
// state
type WatcherLister struct {
	Events <-chan events.Message
	Errors <-chan error
}

// NewWatcherLister is a constructor for the
// WatcherLister type
func NewWatcherLister(ctx context.Context, cli *client.Client) *WatcherLister {
	events, errs := cli.Events(ctx, types.EventsOptions{})

	return &WatcherLister{
		events,
		errs,
	}
}
