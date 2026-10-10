// Copyright (c) 2026 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
)

// getVersion return Psion's version details.
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
func (v *Info) toShortened() string {
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
		Short: "Display the version of tool",
		Run: func(_ *cobra.Command, _ []string) {
			var response string
			var resourceFilesInfo []*ResourceFilesInfo

			versionInfo := getVersion()

			// ignoring error here for now
			resourceFilesInfo, _ = getAllEmbeddedResourceFiles()
			versionInfo.ResourceFiles = resourceFilesInfo

			if shortened {
				response = versionInfo.toShortened()
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
