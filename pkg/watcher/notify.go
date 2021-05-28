package watcher

import (
	"context"
	"time"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/utils"
	"github.com/docker/docker/client"
)

// Notifier sends a diff of the system state
// on the occurence of certain type of events
type Notifier struct {
	Notification      chan core.Diff
	registeredEvenets map[string]bool
}

// NewNotifier registers the events to notify on and
// returns a new Notifier
func NewNotifier(eventsToBeRegistered ...string) *Notifier {
	notifChannel := make(chan core.Diff)

	eventsMap := make(map[string]bool)
	for _, event := range eventsToBeRegistered {
		eventsMap[event] = true
	}

	return &Notifier{
		Notification:      notifChannel,
		registeredEvenets: eventsMap,
	}
}

// Notify creates a diff object and sends it on the Notification channel
func (n *Notifier) Notify(ctx context.Context, cli *client.Client, desired *core.DesiredState, check time.Duration) error {
	feedback := NewFeedback(ctx, cli)
	periodicChecker := time.NewTicker(check)

	for {
		select {
		case event := <-feedback.Events:
			// check if the event recieved is a registered event or not
			if _, eventRegistered := n.registeredEvenets[event.Action]; !eventRegistered {
				continue
			}

			// get the current state of the system
			currentState, err := utils.GetCurrentState(ctx, cli, desired.ContainerType)
			if err != nil {
				return err
			}

			// create the diff object
			diff := core.Diff{
				Current: currentState,
				Desired: desired,
			}

			// send the notification
			n.Notification <- diff
			continue
		// run a periodic check on system state
		case <-periodicChecker.C:
			// get the current state of the system
			currentState, err := utils.GetCurrentState(ctx, cli, desired.ContainerType)
			if err != nil {
				return err
			}

			// create the diff object
			diff := core.Diff{
				Current: currentState,
				Desired: desired,
			}

			// send the notification
			n.Notification <- diff

		case err := <-feedback.Errors:
			return err
		}
	}
}
