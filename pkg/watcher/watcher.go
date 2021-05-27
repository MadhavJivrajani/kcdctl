package watcher

import (
	"context"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
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

// Notifier sends a diff of the system state
// on the occurence of certain type of events
type Notifier struct {
	Notification     chan core.Diff
	RegisteredEvents map[string]bool
}

// NewNotifier registers the events to notify on and
// returns a new Notifier
func NewNotifier(eventsToRegister ...string) *Notifier {
	notifChannel := make(chan core.Diff)

	eventsMap := make(map[string]bool)
	for _, event := range eventsToRegister {
		eventsMap[event] = true
	}

	return &Notifier{
		Notification:     notifChannel,
		RegisteredEvents: eventsMap,
	}
}
