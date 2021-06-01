package declarative

import (
	"context"
	"log"
	"time"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/notifier"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/utils"
	"github.com/docker/docker/client"
)

// set of possible actions for the processor
const (
	// the spawn action indicates to the
	// processor to spawn a new container
	// based on the provided configuration.
	spawn = iota
)

type commandWithSomeOtherStuff struct {
	ctx       *context.Context
	cli       *client.Client
	delta     int
	command   int
	container core.Container
}

func generateCommand(ctx context.Context, cli *client.Client, delta, cmnd int, container core.Container) commandWithSomeOtherStuff {
	return commandWithSomeOtherStuff{
		&ctx,
		cli,
		delta,
		cmnd,
		container,
	}
}

func (cmd commandWithSomeOtherStuff) spawnFunction() error {
	log.Println("Calculated diff:", cmd.delta)
	errChan := make(chan error, cmd.delta)

	for i := 0; i < cmd.delta; i++ {
		go func() {
			log.Println("Spawning container...")
			err := utils.SpawnContainer(*cmd.ctx, cmd.cli, cmd.container)
			if err != nil {
				errChan <- err
				return
			}
			errChan <- nil
			log.Println("Spawn successful...")
		}()
	}

	for i := 0; i < cmd.delta; i++ {
		select {
		case err := <-errChan:
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// The processor executes the command provided
// by the controller.
func processor(command commandWithSomeOtherStuff) error {
	log.Println("System in state drift, attempting reconcile")

	switch command.command {
	case spawn:
		err := command.spawnFunction()
		if err != nil {
			return err
		}
	}

	log.Println("State reconciled")
	return nil
}

// Controller is the controller implementing a control loop and invoking the
// Processor to reconcile the current state and the desired state
func Controller(ctx context.Context, cli *client.Client, eventsToRegister []string, desiredState *core.DesiredState, check time.Duration) error {
	// create a notifier that registers events, on whose
	// occurence, notifications are sent
	notifier := notifier.NewNotifier(eventsToRegister...)

	// start the notification watch
	go notifier.Notify(ctx, cli, desiredState, check)

	// control loop
	for {
		select {
		case <-notifier.Notification:
			// get the current state of the system
			currentState, err := utils.GetCurrentState(ctx, cli, desiredState.ContainerType)
			if err != nil {
				return err
			}

			delta := desiredState.DesiredNum - currentState.CurrentNum
			// TODO: implement excess container deletion
			if delta <= 0 {
				continue
			}

			command := generateCommand(ctx, cli, delta, spawn, desiredState.ContainerType)

			// invoke the processor with the the command that will
			// help reconcil current and desired state.
			err = processor(command)
			if err != nil {
				return err
			}
		}
	}
}
