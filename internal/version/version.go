package version

import (
	"encoding/json"
	"fmt"

	goVersion "go.hein.dev/go-version"
)

// New create a new instance of Version.
func New() *Info {
	return &Info{}
}

// LoadVersion return the Info details.
func (i *Info) LoadVersion(
	version string,
	commit string,
	date string,
) {
	versionOutput := goVersion.New(version, commit, date)

	i.Version = versionOutput.Version
	i.Commit = versionOutput.Commit
	i.Date = versionOutput.Date
}

// ToJSON converts the Info into a JSON String.
func (i *Info) ToJSON() string {
	bytes, _ := json.Marshal(i)

	return string(bytes)
}

// ToShortened converts the Version into a String.
func (i *Info) ToShortened() string {
	return fmt.Sprintf("Version: %s\n", i.Version)
}
