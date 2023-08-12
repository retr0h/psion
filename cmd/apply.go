package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	fileutil "github.com/retr0h/psion/internal/file"
	"github.com/retr0h/psion/pkg/resource/api"
	"github.com/retr0h/psion/pkg/resource/api/v1alpha1"
	"github.com/retr0h/psion/pkg/resource/scheme"
)

var (
	file      string
	resources []api.Manager
)

// applyCmd represents the status command.
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply desired state.",
	Long: `Make consistent with the desired state.
`,
	Run: func(cmd *cobra.Command, args []string) {
		appFs := afero.NewOsFs()

		// todo: loop through list of configs and construct manager resources
		resources := make([]api.Manager, 0, 1)
		fileContent, err := fileutil.Read(appFs, file)
		if err != nil {
			logrus.WithError(err).
				WithField("file", file).
				Fatal("cannot read file")
		}

		s := scheme.New(logger)
		groupVersion := s.Decode(fileContent)

		if groupVersion.Kind == "File" &&
			groupVersion.Version() == "v1alpha1" {
			var resourceKind api.Manager
			resourceKind = v1alpha1.NewFile(
				groupVersion.Version(),
				groupVersion.Kind,
				logger,
				appFs,
				fileContent,
			)
			resources = append(resources, resourceKind)
		}

		for _, resource := range resources {
			resource.Reconcile()
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "todo.")

	// todo: shitty name, change
	rootCmd.MarkPersistentFlagRequired("file")
}
