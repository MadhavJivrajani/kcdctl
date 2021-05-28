package utils

import (
	"context"
	"fmt"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	networkName = "kcd-bangalore-demo"
	ipv4HostIP  = "0.0.0.0"
	ipv6HostIP  = "::"
)

// SpawnContainer creates a new container given the configuration for it
func SpawnContainer(ctx context.Context, cli *client.Client, containterConfig core.Container, netID string) error {
	image := containterConfig.Image

	// pull image
	_, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	// create the container
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port(containterConfig.ContainerPort): []nat.PortBinding{
					{
						HostIP:   ipv4HostIP,
						HostPort: containterConfig.HostPort,
					},
					{
						HostIP:   ipv6HostIP,
						HostPort: containterConfig.HostPort,
					},
				},
			},
			AutoRemove: true,
		},
		nil,
		nil,
		containterConfig.Name,
	)
	if err != nil {
		return err
	}

	// start the container
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	// attach container to the network created on the docker host
	err = cli.NetworkConnect(ctx, netID, resp.ID, nil)
	return err
}

// BootstrapHost bootstraps the docker host with required setup:
// - Creates a new network that the containers will be attached to
// - Creates an nginx container that will act as a load balancer
func BootstrapHost(ctx context.Context, cli *client.Client, lb core.LoadBalancer) (string, error) {
	resp, err := cli.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return "", err
	}

	// create the nginx container that will act as
	// a load balancer
	nginx := core.Container{
		Image:         "nginx",
		Name:          lb.Name,
		ContainerPort: "80/tcp",
		HostPort:      lb.ExposedPort,
	}
	err = SpawnContainer(ctx, cli, nginx, resp.ID)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// GetCurrentState gets the current state of the system based on the common name prefix.
func GetCurrentState(ctx context.Context, cli *client.Client, containterConfig core.Container) (*core.CurrentState, error) {
	// containterConfig.Name is the prefix of
	// all container names that run as part of
	// this 'service', to filter based on these
	// prefixes, add a filter with teh regex
	// ^name*
	listOpts := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "name",
				Value: fmt.Sprintf("^%s*", containterConfig.Name),
			},
		),
	}
	ctrs, err := cli.ContainerList(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	currentState := &core.CurrentState{
		CurrentNum:    len(ctrs),
		ContainerType: containterConfig,
	}
	return currentState, nil
}
