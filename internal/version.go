package internal

// VersionManager manager responsible for Version operations.
type VersionManager interface {
	LoadVersion(
		version string,
		commit string,
		date string,
	)
	ToJSON() string
	ToShortened() string
}
