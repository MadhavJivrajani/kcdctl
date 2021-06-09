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
				"serviceSelector=configName",
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

	hostConfig := &container.HostConfig{
		AutoRemove: true,
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
