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

package imperative

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/MadhavJivrajani/kcdctl/pkg/core"
	"github.com/MadhavJivrajani/kcdctl/pkg/notifier"
	"github.com/MadhavJivrajani/kcdctl/pkg/utils"
	"github.com/docker/docker/client"
)

// StartObserving starts a notification channel that recievs registered events and also periodic
// events to check state.
func StartObserving(ctx context.Context, cli *client.Client, eventsToRegister []string, desired *core.DesiredState, check time.Duration, wg *sync.WaitGroup) error {
	defer wg.Done()
	// create a notifier
	notifier := notifier.NewNotifier(eventsToRegister...)

	// start the notification watch
	go notifier.Notify(ctx, cli, desired, check)

	var reconciled bool
	for {
		event := <-notifier.Notification
		if event != "check" {
			log.Println("Event recieved:", event)
		}
		currentState, err := utils.GetCurrentState(ctx, cli, desired.ContainerType)
		if err != nil {
			return err
		}
		diff := desired.DesiredNum - currentState.CurrentNum
		if diff != 0 {
			log.Println("Current diff in terms of replicas:", diff)
			log.Println("Awaiting reconcilitation...")
		} else {
			if !reconciled {
				reconciled = !reconciled
				log.Println("State reconciled")
			}
		}
	}
}
