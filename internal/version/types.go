package version

// Info creates a formattable struct for output.
type Info struct {
	Version       string               `json:"version,omitempty"`
	Commit        string               `json:"commit,omitempty"`
	Date          string               `json:"date,omitempty"`
	ResourceFiles []*ResourceFilesInfo `json:"resource_files,omitempty"`
}

// ResourceFilesInfo describes embedded files.
type ResourceFilesInfo struct {
	Path     string `json:"path,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Type     string `json:"type,omitempty"`
}
