package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/retr0h/psion/internal/config"
	fileutil "github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
	"github.com/retr0h/psion/pkg/resource/api/v1alpha1"
)

var file string

// resources []api.Manager

// applyCmd represents the status command.
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply desired state.",
	Long: `Make consistent with the desired state.
`,
	Run: func(cmd *cobra.Command, args []string) {
		appFs := afero.NewOsFs()
		plan := false

		// todo: loop through list of configs and construct manager resources
		resources := make([]api.Manager, 0, 1)
		fileContent, err := fileutil.Read(appFs, file)
		if err != nil {
			logrus.WithError(err).
				WithField("file", file).
				Fatal("cannot read file")
		}

		runtimeConfig, err := config.LoadRuntimeConfig(fileContent)
		if err != nil {
			logrus.WithError(err).
				WithField("file", file).
				Fatal("cannot load file")
		}

		if runtimeConfig.APIVersion != v1alpha1.FileAPIVersion {
			logrus.WithField("apiVersion", runtimeConfig.APIVersion).
				Fatal("invalid api version")
		}

		if runtimeConfig.Kind != v1alpha1.FileKind {
			logrus.WithField("kind", runtimeConfig.Kind).
				Fatal("invalid kind")
		}

		var resourceKind api.Manager = v1alpha1.NewFile(
			logger,
			appFs,
			plan,
		)

		if err := yaml.Unmarshal(fileContent, resourceKind); err != nil {
			logrus.WithError(err).
				WithField("file", file).
				Fatal("cannot unmarshal file")
		}
		resources = append(resources, resourceKind)

		for _, resource := range resources {
			if err := resource.Reconcile(); err != nil {
        logrus.WithError(err).
          Fatal("cannot reconcile")
		  }
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "todo.")

	// todo: shitty name, change
	_ = rootCmd.MarkPersistentFlagRequired("file")
}
