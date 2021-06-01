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

package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const networkName = "kcd-bangalore-demo"

// cleanupCmd represents the cleanup command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "cleanup removes all containers that were created as part of a given configuration file",

	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		if file == "" {
			return fmt.Errorf("no config file passed! please pass one :(")
		}
		config, err := utils.ReadConfig(file)
		if err != nil {
			return err
		}

		err = cleanup(config)
		return err
	},
}

func init() {
	rootCmd.AddCommand(cleanupCmd)

	cleanupCmd.Flags().StringP(
		"file",
		"f",
		"",
		"file is the yaml configuarion file using which containers will be cleaned up",
	)
}

func stopContainersByLabel(ctx context.Context, cli *client.Client, label string) error {
	listOpts := types.ContainerListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "label",
				Value: fmt.Sprintf("configName=%s", label),
			},
		),
	}

	containers, err := cli.ContainerList(ctx, listOpts)
	if err != nil {
		return err
	}

	errChan := make(chan error, len(containers))

	// stop containers
	for _, ctr := range containers {
		go func(ctr types.Container) {
			err := cli.ContainerStop(ctx, ctr.ID, nil)
			errChan <- err
		}(ctr)
	}

	for i := 0; i < len(containers); i++ {
		select {
		case err := <-errChan:
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func removeNetworkByLabel(ctx context.Context, cli *client.Client, label string) error {
	listOpts := types.NetworkListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "name",
				Value: label,
			},
		),
	}

	networks, err := cli.NetworkList(ctx, listOpts)
	if err != nil {
		return err
	}
	for _, network := range networks {
		err := cli.NetworkRemove(ctx, network.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func cleanup(config utils.Config) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	ctx := context.Background()

	// remove containers running the application
	log.Println("Cleaning up containers running application...")
	err = stopContainersByLabel(ctx, cli, config.Spec.Template.Name)
	if err != nil {
		return err
	}

	// remove the load balancer
	log.Println("Cleaning up load balancer container...")
	err = stopContainersByLabel(ctx, cli, config.LoadBalancer.Name)
	if err != nil {
		return err
	}

	// cleaning up network
	log.Println("Cleaning up network created...")
	err = removeNetworkByLabel(ctx, cli, networkName)

	return err
}
