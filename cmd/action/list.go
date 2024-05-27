//
// Copyright 2024 Stacklok, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package action

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/stacklok/frizbee/internal/cli"
	"github.com/stacklok/frizbee/pkg/config"
	"github.com/stacklok/frizbee/pkg/replacer"
	"os"
	"strconv"
)

// CmdList represents the one sub-command
func CmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists the used github actions",
		Long: `This utility lists all the github actions used in the workflows

Example: 
	frizbee action list -d .github/workflows
`,
		Aliases:      []string{"ls"},
		RunE:         list,
		SilenceUsage: true,
	}

	cli.DeclareFrizbeeFlags(cmd, ".github/workflows")
	cmd.Flags().StringP("output", "o", "table", "output format. Can be 'json' or 'table'")

	return cmd
}

func list(cmd *cobra.Command, _ []string) error {
	// Extract the CLI flags from the cobra command
	cliFlags, err := cli.NewHelper(cmd)
	if err != nil {
		return err
	}

	// Set up the config
	cfg, err := config.FromCommand(cmd)
	if err != nil {
		return err
	}

	// Create a new replacer
	r := replacer.New(cfg).
		WithUserRegex(cliFlags.Regex).
		WithGitHubClient(os.Getenv(cli.GitHubTokenEnvKey))

	// List the references in the directory
	res, err := r.ListGitHibActions(cliFlags.Dir)
	if err != nil {
		return err
	}

	output := cmd.Flag("output").Value.String()
	switch output {
	case "json":
		jsonBytes, err := json.MarshalIndent(res.Entities, "", "  ")
		if err != nil {
			return err
		}
		jsonString := string(jsonBytes)

		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", jsonString)
		return nil
	case "table":
		table := tablewriter.NewWriter(cmd.OutOrStdout())
		table.SetHeader([]string{"No", "Type", "Name", "Ref"})
		for i, a := range res.Entities {
			table.Append([]string{strconv.Itoa(i + 1), a.Type, a.Name, a.Ref})
		}
		table.Render()
		return nil
	default:
		return fmt.Errorf("unknown output format: %s", output)
	}
}
