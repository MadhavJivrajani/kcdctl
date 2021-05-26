package main

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	eventOpts := types.EventsOptions{}
	events, errChan := cli.Events(context.Background(), eventOpts)

	for {
		select {
		case event := <-events:
			fmt.Println(event.Action)
		case err := <-errChan:
			panic(err)
		}
	}
}
