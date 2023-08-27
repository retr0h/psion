package cmd

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/retr0h/psion/internal/config"
	fileutil "github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
	"github.com/retr0h/psion/pkg/resource/api/v1alpha1"
)

var file string

// applyCmd represents the status command.
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply desired state.",
	Long: `Make consistent with the desired state.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// By the time we reach this point, we know that the arguments were
		// properly parsed, and we don't want to show the usage if an error
		// occurs.
		cmd.SilenceUsage = true

		appFs := afero.NewOsFs()
		plan := false

		// todo: loop through list of configs and construct manager resources
		resources := make([]api.Manager, 0, 1)
		fileContent, err := fileutil.Read(appFs, file)
		if err != nil {
			return fmt.Errorf("cannot read file: %w", err)
		}

		runtimeConfig, err := config.LoadRuntimeConfig(fileContent)
		if err != nil {
			return fmt.Errorf("cannot load file: %w", err)
		}

		if runtimeConfig.APIVersion != v1alpha1.FileAPIVersion {
			return fmt.Errorf("invalid apiVersion: %s", runtimeConfig.APIVersion)
		}

		if runtimeConfig.Kind != v1alpha1.FileKind {
			return fmt.Errorf("invalid kind: %s", runtimeConfig.Kind)
		}

		var resourceKind api.Manager = v1alpha1.NewFile(
			logger,
			appFs,
			plan,
		)

		if err := yaml.Unmarshal(fileContent, resourceKind); err != nil {
			return fmt.Errorf("cannot unmarshal file: %w", err)
		}
		resources = append(resources, resourceKind)

		for _, resource := range resources {
			if err := resource.Reconcile(); err != nil {
				return fmt.Errorf("cannot reconcile: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "todo.")

	// todo: shitty name, change
	_ = applyCmd.MarkPersistentFlagRequired("file")
}
