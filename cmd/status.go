package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
)

// statusCmd represents the plan command.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the current status",
	Long: `Display the current status of the state file.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		style := newCLIStyle(os.Stdout)
		fmt.Printf("\n  %s\n\n", style.title.Render("RESOURCE STATUS"))

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
			style.title.Render("Name"),
			style.title.Render("Status"),
			style.title.Render("Kind"),
			style.title.Render("APIVersion"),
			style.title.Render("Conditions"),
		})
		for _, resource := range state.GetItems() {
			tConditions := generateInnerTable()
			for _, condition := range resource.Status.Conditions {
				tConditions.AppendRow(table.Row{"Type", condition.Type})
				tConditions.AppendRow(table.Row{"Status", style.phase(string(condition.Status))})
				tConditions.AppendRow(table.Row{"Message", condition.Message})
				tConditions.AppendRow(table.Row{"Reason", condition.Reason})
				tConditions.AppendRow(table.Row{"Got", condition.Got})
				tConditions.AppendRow(table.Row{"Want", condition.Want})
				tConditions.AppendSeparator()
			}
			t.AppendRow(
				table.Row{
					style.accent.Render(resource.Name),
					style.phase(string(resource.Phase)),
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
				style.title.Render("Status"),
				style.phase(state.GetStatusString()),
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
