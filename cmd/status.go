// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/cli"
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

// statusCmd represents the plan command.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the current status",
	Long: `Display the current status of the state file.
`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		// By the time we reach this point, we know that the arguments were
		// properly parsed, and we don't want to show the usage if an error
		// occurs.
		cmd.SilenceUsage = true

		var (
			fileManager internal.FileManager = file.New(appFs)
			state       api.StateManager     = api.NewState(fileManager, stateFile)
		)

		state, err := state.GetState()
		if err != nil {
			return fmt.Errorf("cannot get state: %w", err)
		}

		style := cli.NewTheme(os.Stdout, color)
		fmt.Printf("\n  %s\n\n", style.Title("RESOURCE STATUS"))

		generateInnerTable := func() table.Writer {
			tw := table.NewWriter()
			tw.Style().Options.DrawBorder = false
			tw.Style().Options.SeparateColumns = false

			return tw
		}

		t := table.NewWriter()
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateColumns = false
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{
			style.Title("Name"),
			style.Title("Status"),
			style.Title("Kind"),
			style.Title("APIVersion"),
			style.Title("Conditions"),
		})
		for _, resource := range state.GetItems() {
			tConditions := generateInnerTable()
			for _, condition := range resource.Status.Conditions {
				tConditions.AppendRow(table.Row{"Type", condition.Type})
				tConditions.AppendRow(table.Row{"Status", style.Phase(string(condition.Status))})
				tConditions.AppendRow(table.Row{"Message", condition.Message})
				tConditions.AppendRow(table.Row{"Reason", condition.Reason})
				tConditions.AppendRow(table.Row{"Got", condition.Got})
				tConditions.AppendRow(table.Row{"Want", condition.Want})
				tConditions.AppendSeparator()
			}
			t.AppendRow(
				table.Row{
					style.Accent.Render(resource.Name),
					style.Phase(string(resource.Phase)),
					resource.Kind,
					resource.APIVersion,
					tConditions.Render(),
				},
			)
			t.AppendSeparator()
		}
		t.AppendSeparator()
		t.AppendFooter(
			table.Row{
				style.Title("Status"),
				style.Phase(state.GetStatusString()),
				"",
				"",
				"",
			},
		)
		t.Render()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
