package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/retr0h/psion/pkg/resource/api"
)

// applyCmd represents the apply command.
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply desired state",
	Long: `Make consistent with the desired state.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// By the time we reach this point, we know that the arguments were
		// properly parsed, and we don't want to show the usage if an error
		// occurs.
		cmd.SilenceUsage = true

		plan := false
		resources, err := loadAllEmbeddedResourceFiles(plan)
		if err != nil {
			return fmt.Errorf("cannot walk dir: %w", err)
		}

		stateResources := make([]*api.Resource, 0, 1)
		for _, resource := range resources {
			if err := resource.Reconcile(); err != nil {
				return fmt.Errorf("cannot reconcile: %w", err)
			}

			state := resource.GetState()
			stateResources = append(stateResources, state)
		}

		if err := writeStateFile(stateResources); err != nil {
			return fmt.Errorf("cannot write state file: %w", err)
		}

		logger.Info(
			"wrote state file",
			slog.String("StateFile", stateFile),
		)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
