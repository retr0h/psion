package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	// "github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
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

		state, err := getState()
		if err != nil {
			return fmt.Errorf("cannot get state: %w", err)
		}

		generateInnerTable := func() table.Writer {
			tw := table.NewWriter()
			tw.Style().Options.DrawBorder = false

			return tw
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name", "Status", "Kind", "APIVersion", "Conditions"})
		for _, resource := range state.Items {
			tConditions := generateInnerTable()
			for _, condition := range resource.Status.Conditions {
				tConditions.AppendRow(table.Row{"Type", condition.Type})
				tConditions.AppendRow(table.Row{"Status", condition.Status})
				tConditions.AppendRow(table.Row{"Message", condition.Message})
				tConditions.AppendRow(table.Row{"Reason", condition.Reason})
				tConditions.AppendRow(table.Row{"Got", condition.Got})
				tConditions.AppendRow(table.Row{"Want", condition.Want})
				tConditions.AppendSeparator()
			}
			t.AppendRow(
				table.Row{
					resource.Name,
					resource.Phase,
					resource.Kind,
					resource.APIVersion,
					tConditions.Render(),
				},
			)
			t.AppendSeparator()
		}
		t.AppendSeparator()
		t.Render()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
