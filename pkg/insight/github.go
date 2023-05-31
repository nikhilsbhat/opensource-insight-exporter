package insight

// GitRelease holds information of the releases downloaded metrics of GitHub.
type GitRelease struct {
	TagName string   `json:"tag_name,omitempty" yaml:"tag_name,omitempty"`
	Name    string   `json:"name,omitempty" yaml:"name,omitempty"`
	Assets  []Assets `json:"assets,omitempty" yaml:"assets,omitempty"`
}

// Assets holds download count information of a specific release of GitHub.
type Assets struct {
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	DownloadCount int    `json:"download_count,omitempty" yaml:"download_count,omitempty"`
}
