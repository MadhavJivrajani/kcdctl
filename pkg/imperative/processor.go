package imperative

import (
	"context"
	"fmt"
	"log"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/utils"
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

func Processor(cmd toExecute) error {
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
