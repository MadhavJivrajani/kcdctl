package main

import (
	"context"
	"log"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/controller"
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

	log.Println("Bootstrapping host...")
	netID, err := utils.BootstrapHost(ctx, cli, core.LoadBalancer{
		Name:        "lb",
		ExposedPort: "9090",
		TargetPort:  "8080",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Bootstrap successful")

	container := core.Container{
		Name:          "test",
		Image:         "nginx",
		HostPort:      "8080",
		ContainerPort: "80",
		Network:       netID,
	}
	events := []string{"kill", "stop", "die", "destroy"}
	desiredState := &core.DesiredState{
		DesiredNum:    2,
		ContainerType: container,
	}

	log.Println("Desired state:", desiredState.DesiredNum)
	log.Println("Starting controller...")
	err = controller.Controller(ctx, cli, events, desiredState)
	if err != nil {
		log.Fatal(err)
	}
}
