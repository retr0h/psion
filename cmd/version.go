package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

// versionCmd represents the version command
var (
	shortened  = false
	version    = "dev"
	commit     = "none"
	date       = "unknown"
	output     = "json"
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display the version of geoip-updater",
		Run: func(cmd *cobra.Command, args []string) {
			resp := goVersion.FuncWithOutput(shortened, version, commit, date, output)

			fmt.Println(resp)
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&shortened, "short", "s", false, "Print just the version number.")
}
