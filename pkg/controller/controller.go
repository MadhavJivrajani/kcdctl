package controller

import (
	"context"
	"log"
	"time"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/utils"
	"github.com/docker/docker/client"
)

// Processor decides an action to take place based on the current state
// of the system and the desired state of the system
func Processor(ctx context.Context, cli *client.Client, currentState *core.CurrentState, desiredState *core.DesiredState) error {
	currentNum := currentState.CurrentNum
	desiredNum := desiredState.DesiredNum

	delta := desiredNum - currentNum
	if delta == 0 {
		return nil
	}
	log.Println("system in state drift, attempting reconcile")
	log.Println("current state:", currentState.CurrentNum)
	for i := 0; i < delta; i++ {
		// reconcile state
		log.Println("spawning container...")
		err := utils.SpawnContainer(ctx, cli, desiredState.ContainerType)
		if err != nil {
			return err
		}
		log.Println("spawn successful")
	}
	log.Println("state reconciled")
	return nil
}

// Controller is the controller implementing a control loop and invoking the
// Processor to reconcile the current state and the desired state
func Controller(ctx context.Context, cli *client.Client, eventsToRegister []string, desiredState *core.DesiredState, check time.Duration) error {
	// control loop
	for {
		currentState, err := utils.GetCurrentState(ctx, cli, desiredState.ContainerType)
		if err != nil {
			return err
		}
		// invoke the Processor to try and reconcile current and
		// desired state
		err = Processor(ctx, cli, currentState, desiredState)
		if err != nil {
			return err
		}
		time.Sleep(check)
	}
}
