package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

// getVersion return the sensor's VersionInfo details.
func getVersion() *Info {
	versionOutput := goVersion.New(version, commit, date)

	output := &Info{
		Version: versionOutput.Version,
		Commit:  versionOutput.Commit,
		Date:    versionOutput.Date,
	}

	return output
}

// toJSON converts the Info into a JSON String.
func (v *Info) toJSON() string {
	bytes, _ := json.Marshal(v)

	return string(bytes)
}

// toShortened converts the Version into a String.
func (v *Info) ToShortened() string {
	return fmt.Sprintf("Version: %s\n", v.Version)
}

// versionCmd represents the version command.
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
			var response string
			var resourceFilesInfo []*ResourceFilesInfo
			var err error

			versionInfo := getVersion()

			resourceFilesInfo, err = getAllEmbeddedResourceFiles()
			if err != nil {
				resourceFilesInfo = resourceFilesInfo
			}

			versionInfo.ResourceFiles = resourceFilesInfo

			if shortened {
				response = versionInfo.ToShortened()
			} else {
				response = versionInfo.toJSON()
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
