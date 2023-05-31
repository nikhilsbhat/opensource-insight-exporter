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
