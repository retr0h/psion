package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/retr0h/psion/internal"
	"github.com/retr0h/psion/internal/version"
)

// versionCmd represents the version command.
var (
	shortened  = false
	ver        = "dev"
	commit     = "none"
	date       = "unknown"
	output     = "json"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display the version of tool",
		Run: func(cmd *cobra.Command, args []string) {
			var vm internal.VersionManager = version.New()
			var response string
			// var resourceFilesInfo []*ResourceFilesInfo

			vm.LoadVersion(ver, commit, date)

			// // ignoring error here for now
			// resourceFilesInfo, _ = getAllEmbeddedResourceFiles()
			// versionInfo.ResourceFiles = resourceFilesInfo

			if shortened {
				response = vm.ToShortened()
			} else {
				response = vm.ToJSON()
			}

			fmt.Printf("%s\n", response)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
	versionCmd.Flags().
		StringVarP(&output, "output", "o", "json", "Output format. One of 'yaml' or 'json'.")
}
