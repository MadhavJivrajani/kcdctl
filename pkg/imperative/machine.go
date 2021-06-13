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

package imperative

import (
	"context"
	"fmt"
	"log"

	"github.com/MadhavJivrajani/kcdctl/pkg/core"
	"github.com/MadhavJivrajani/kcdctl/pkg/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

const networkName = "kcd-bangalore-demo"

type Command struct {
	Config utils.Config
	Diff   int
}

func (cmd Command) Spawn() error {
	return spawnDiff(cmd.Config, cmd.Diff)
}

type toExecute func() error

func Machine(cmd toExecute) error {
	err := cmd()
	return err
}

func spawnDiff(config utils.Config, diff int) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	netID, err := getNetworkID(ctx, cli)
	if err != nil {
		return err
	}

	log.Println("Recieved diff:", diff)
	// form the container type based on the config
	container := core.Container{
		Name:    config.Spec.Template.Name,
		Image:   config.Spec.Template.Image,
		Network: netID,
	}
	errChan := make(chan error, diff)

	// spawn diff containers
	for i := 0; i < diff; i++ {
		go func() {
			log.Println("Spawning container...")
			err := utils.SpawnContainer(ctx, cli, container)
			if err != nil {
				errChan <- err
				return
			}
			errChan <- nil
			log.Println("Spawn successful...")
		}()
	}

	for i := 0; i < diff; i++ {
		err := <-errChan
		if err != nil {
			return err
		}
	}
	return nil
}

func getNetworkID(ctx context.Context, cli *client.Client) (string, error) {
	listOpts := types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "name",
				Value: networkName,
			},
		),
	}

	networks, err := cli.NetworkList(ctx, listOpts)
	if err != nil {
		return "", err
	}
	if len(networks) == 0 {
		return "", fmt.Errorf("no networks with name %s exist", networkName)
	}
	return networks[0].ID, nil
}
