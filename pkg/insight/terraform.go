package insight

// ProviderDownloadSummary holds the information of the download metrics of a Terraform provider.
type ProviderDownloadSummary struct {
	Data Data `json:"data,omitempty" yaml:"data,omitempty"`
}

// Data holds the download data of a Terraform provider.
type Data struct {
	Type       string     `json:"type,omitempty" yaml:"type,omitempty"`
	ID         string     `json:"id,omitempty" yaml:"id,omitempty"`
	Attributes Attributes `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}

// Attributes holds the download attributes of a Terraform provider.
type Attributes struct {
	Month int `json:"month,omitempty" yaml:"month,omitempty"`
	Total int `json:"total,omitempty" yaml:"total,omitempty"`
	Week  int `json:"week,omitempty" yaml:"week,omitempty"`
	Year  int `json:"year,omitempty" yaml:"year,omitempty"`
}

// ProviderVersion Holds the information of all available versions of a provider.
type ProviderVersion struct {
	ID        string   `json:"id,omitempty" yaml:"id,omitempty"`
	Versions  []string `json:"versions,omitempty" yaml:"versions,omitempty"`
	Source    string   `json:"source,omitempty" yaml:"source,omitempty"`
	Owner     string   `json:"owner,omitempty" yaml:"owner,omitempty"`
	NameSpace string   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

// Version Holds the version of a provider.
type Version struct {
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
}
