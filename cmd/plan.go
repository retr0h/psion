package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// planCmd represents the plan command.
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Preview the changes to be made",
	Long: `Plan the changes to make consistent with the desired state.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// By the time we reach this point, we know that the arguments were
		// properly parsed, and we don't want to show the usage if an error
		// occurs.
		cmd.SilenceUsage = true

		plan := true
		resources, err := loadAllEmbeddedResourceFiles(plan)
		if err != nil {
			return fmt.Errorf("cannot walk dir: %w", err)
		}

		for _, resource := range resources {
			if err := resource.Reconcile(); err != nil {
				return fmt.Errorf("cannot reconcile: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
