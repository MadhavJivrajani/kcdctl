package main

import (
	"context"
	"log"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/watcher"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	n := watcher.NewNotifier("kill", "stop")
	ctr := core.Container{}
	desired := &core.DesiredState{
		DesiredNum:    2,
		ContainerType: ctr,
	}

	go n.Notify(ctx, cli, desired)

	for {
		select {
		case diff := <-n.Notification:
			log.Println(*diff.Current, *diff.Desired)
		}
	}
}
