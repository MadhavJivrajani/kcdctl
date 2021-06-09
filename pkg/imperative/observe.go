package imperative

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/notifier"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/utils"
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
