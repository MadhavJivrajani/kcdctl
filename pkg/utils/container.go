package utils

import (
	"context"
	"fmt"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	networkName = "kcd-bangalore-demo"
	ipv4HostIP  = "0.0.0.0"
	ipv6HostIP  = "::"
)

func spawnLoadBalancer(ctx context.Context, cli *client.Client, config Config, netID string) error {
	image := config.LoadBalancer.Image

	// pull image
	_, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%s/tcp", config.LoadBalancer.ContainerPort)): []nat.PortBinding{
				{
					HostIP:   ipv4HostIP,
					HostPort: config.LoadBalancer.ExposedPort,
				},
				{
					HostIP:   ipv6HostIP,
					HostPort: config.LoadBalancer.ExposedPort,
				},
			},
		},
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:     mount.TypeBind,
				Source:   "/var/run/docker.sock",
				Target:   "/var/run/docker.sock",
				ReadOnly: true,
			},
		},
	}

	// create the container
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Labels: map[string]string{
				"configName": config.LoadBalancer.Name,
			},
			Env: []string{
				fmt.Sprintf("containerPort=%s", config.LoadBalancer.ContainerPort),
				fmt.Sprintf("targetPort=%s", config.LoadBalancer.TargetPort),
				fmt.Sprintf("serviceSelector=configName"),
				fmt.Sprintf("serviceSelectorValue=%s", config.Spec.Template.Name),
			},
			ExposedPorts: nat.PortSet{
				nat.Port(fmt.Sprintf("%s/tcp", config.LoadBalancer.ContainerPort)): struct{}{},
			},
		},
		hostConfig,
		nil,
		nil,
		"",
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

// SpawnContainer creates a new container given the configuration for it
func SpawnContainer(ctx context.Context, cli *client.Client, containterConfig core.Container) error {
	image := containterConfig.Image

	// pull image
	_, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	var hostConfig *container.HostConfig
	if containterConfig.HostPort != "" {
		hostConfig = &container.HostConfig{
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
		}
	}

	// create the container
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Labels: map[string]string{
				"configName": containterConfig.Name,
			},
		},
		hostConfig,
		nil,
		nil,
		"",
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
	err = cli.NetworkConnect(ctx, containterConfig.Network, resp.ID, nil)
	return err
}

// BootstrapHost bootstraps the docker host with required setup:
// - Creates a new network that the containers will be attached to
// - Creates an nginx container that will act as a load balancer
func BootstrapHost(ctx context.Context, cli *client.Client, config Config) (string, error) {
	resp, err := cli.NetworkCreate(ctx, networkName, types.NetworkCreate{})
	if err != nil {
		return "", err
	}

	// create the load balancer
	err = spawnLoadBalancer(ctx, cli, config, resp.ID)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

// GetCurrentState gets the current state of the system based on the common name prefix.
func GetCurrentState(ctx context.Context, cli *client.Client, containterConfig core.Container) (*core.CurrentState, error) {
	ctrs, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	currentState := &core.CurrentState{
		// len(ctrs) - 1 is to take
		// care of the fact that one
		// of the containers is the
		// load balancer
		CurrentNum:    len(ctrs) - 1,
		ContainerType: containterConfig,
	}

	return currentState, nil
}
