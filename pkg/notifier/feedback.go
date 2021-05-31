package notifier

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

// feedback watches the state of the system
// and lists events generated on change in
// system state
type feedback struct {
	Events <-chan events.Message
	Errors <-chan error
}

// newFeedback is a constructor for the feedback type
func newFeedback(ctx context.Context, cli *client.Client) *feedback {
	events, errs := cli.Events(ctx, types.EventsOptions{})

	return &feedback{
		events,
		errs,
	}
}
