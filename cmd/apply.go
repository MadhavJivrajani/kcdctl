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

	"github.com/MadhavJivrajani/kcd-bangalore/pkg/controller"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/core"
	"github.com/MadhavJivrajani/kcd-bangalore/pkg/utils"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const (
	declarativeMode = "declarative"
	imperativeMode  = "imperative"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "applies a configuarion to the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			return err
		}
		// validate mode passed
		if mode != declarativeMode && mode != imperativeMode {
			return fmt.Errorf("invalid mode string passed")
		}

		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}
		// validate a config file is passed because UX is not existent
		if file == "" {
			return fmt.Errorf("no config file passed! please pass one :(")
		}

		config, err := utils.ReadConfig(file)
		if err != nil {
			return err
		}

		if mode == declarativeMode {
			err = startDeclarativeSystem(config)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringP(
		"mode",
		"m",
		declarativeMode,
		"mode is the mode the system should be run in: imperative/declarative",
	)

	applyCmd.Flags().StringP(
		"file",
		"f",
		"",
		"file is the yaml configuarion file that the user wants to apply on the system",
	)
}

func startDeclarativeSystem(config utils.Config) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	ctx := context.Background()

	// bootstrapping host
	log.Println("Bootstrapping host...")
	netID, err := utils.BootstrapHost(ctx, cli, config)
	if err != nil {
		return err
	}
	log.Println("Bootstrap successful")

	// form the container type based on the config
	container := core.Container{
		Name:    config.Spec.Template.Name,
		Image:   config.Spec.Template.Image,
		Network: netID,
	}

	// events to potentially respond to
	events := []string{"kill", "stop", "die", "destroy"}

	// desired state that the user wants the system
	// to be in
	desiredState := &core.DesiredState{
		DesiredNum:    config.Spec.Replicas,
		ContainerType: container,
	}

	log.Println("Desired number of replicas:", desiredState.DesiredNum)
	log.Println("Starting controller...")
	// start the controller
	err = controller.Controller(ctx, cli, events, desiredState)

	return err
}
