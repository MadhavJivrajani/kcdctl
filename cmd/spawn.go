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
	"fmt"

	"github.com/MadhavJivrajani/kcdctl/pkg/imperative"
	"github.com/MadhavJivrajani/kcdctl/pkg/utils"
	"github.com/spf13/cobra"
)

// spawnCmd represents the spawn command
var spawnCmd = &cobra.Command{
	Use:   "spawn",
	Short: "spawn creates containers",
	RunE: func(cmd *cobra.Command, args []string) error {
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

		diff, err := cmd.Flags().GetInt("diff")
		if err != nil {
			return err
		}

		command := imperative.Command{
			Config: config,
			Diff:   diff,
		}

		err = imperative.Processor(command.Spawn)
		return err
	},
}

func init() {
	rootCmd.AddCommand(spawnCmd)

	spawnCmd.Flags().IntP(
		"diff",
		"d",
		1,
		"diff is the number of containers to spawn (right now, assumed to be +ve)",
	)

	spawnCmd.Flags().StringP(
		"file",
		"f",
		"",
		"file is the yaml configuarion file that the user wants to apply on the system",
	)
}
