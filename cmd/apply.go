package cmd

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/internal/state"
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

		var (
			fileManager internal.FileManager = file.New(appFs)
			state       StateManager         = state.New(fileManager, stateFile)
		)

		plan := false
		resources, err := loadAllEmbeddedResourceFiles(plan)
		if err != nil {
			return fmt.Errorf("cannot walk dir: %w", err)
		}

		for _, resource := range resources {
			if err := resource.Reconcile(); err != nil {
				return fmt.Errorf("cannot reconcile: %w", err)
			}
			state.SetItems(resource.GetState())
		}

		if err := state.SetState(); err != nil {
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
