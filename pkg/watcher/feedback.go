package watcher

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

// Feedback watches the state of the system
// and lists events generated on change in
// system state
type Feedback struct {
	Events <-chan events.Message
	Errors <-chan error
}

// NewFeedback is a constructor for the Feedback type
func NewFeedback(ctx context.Context, cli *client.Client) *Feedback {
	events, errs := cli.Events(ctx, types.EventsOptions{})

	return &Feedback{
		events,
		errs,
	}
}
