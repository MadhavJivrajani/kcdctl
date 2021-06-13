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

package declarative

import (
	"context"
	"log"
	"time"

	"github.com/MadhavJivrajani/kcdctl/pkg/core"
	"github.com/MadhavJivrajani/kcdctl/pkg/notifier"
	"github.com/MadhavJivrajani/kcdctl/pkg/utils"
	"github.com/docker/docker/client"
)

// set of possible actions for the processor
const (
	// the spawn action indicates to the
	// processor to spawn a new container
	// based on the provided configuration.
	spawn = iota
)

// command represents a command that the controller
// passes to the processor to execute.
type command func() error

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
		err := <-errChan
		if err != nil {
			return err
		}
	}
	return nil
}

// The processor executes the command provided
// by the controller.
func processor(cmd command) error {
	log.Println("System in state drift, attempting reconcile")

	err := cmd()
	if err != nil {
		return err
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
		<-notifier.Notification
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
		err = processor(command.spawnFunction)
		if err != nil {
			return err
		}
	}
}
