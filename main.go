package main

import (
	"context"
	"log"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/utils"
	"github.com/docker/docker/client"
)

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	netID, err := utils.BootstrapHost(ctx, cli)
	if err != nil {
		log.Fatal(err)
	}

	container := core.Container{
		Image:         "nginx",
		Name:          "test-container",
		HostPort:      "8080",
		ContainerPort: "80/tcp",
	}
	err = utils.SpawnContainer(ctx, cli, container, netID)
	if err != nil {
		log.Fatal(err)
	}
}
