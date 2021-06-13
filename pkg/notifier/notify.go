/*
Copyright © 2021 Madhav Jivrajani madhav.jiv@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package notifier

import (
	"context"
	"time"

	"github.com/MadhavJivrajani/kcdctl/pkg/core"
	"github.com/docker/docker/client"
)

const checkEvent = "check"

// Notifier sends a notification which is
// in the form of the event that occured
type Notifier struct {
	Notification      chan string
	registeredEvenets map[string]bool
}

// NewNotifier registers the events to notify on and
// returns a new Notifier
func NewNotifier(eventsToBeRegistered ...string) *Notifier {
	notifChannel := make(chan string)

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
	// Create a periodic check in case events either miss or the
	// state checked in the controller does not accurately represent
	// the change in state of the system.
	periodicChecker := time.NewTicker(check)

	feedback := newFeedback(ctx, cli)
	for {
		select {
		case event := <-feedback.Events:
			// check if the event recieved is a registered event or not
			if _, eventRegistered := n.registeredEvenets[event.Action]; !eventRegistered {
				continue
			}

			n.Notification <- event.Action
		case err := <-feedback.Errors:
			return err
		case <-periodicChecker.C:
			n.Notification <- checkEvent
		}
	}
}
