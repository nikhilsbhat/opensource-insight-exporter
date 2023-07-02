package insight

import (
	"os"

	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/common"
)

// Insight holds the download metrics for a particular project.
type Insight struct {
	Platform string `json:"platform,omitempty" yaml:"platform,omitempty"`
	ID       string `json:"id,omitempty" yaml:"id,omitempty"`
	Summary  any    `json:"summary,omitempty" yaml:"summary,omitempty"`
}

// GetInsight gets insight into the downloads of all specified projects.
func (cfg *Config) GetInsight() []Insight {
	insight := make([]Insight, 0)
	for _, source := range cfg.Sources {
		source.SetLogger(cfg.logger)

		caContent := make([]byte, 0)
		if len(source.CaFilePath) != 0 {
			ca, err := os.ReadFile(source.CaFilePath)
			if err != nil {
				cfg.logger.Errorf("reading ca file errored with %v", err)
			}
			caContent = ca
		}

		httpClient := source.NewClient(caContent)

		for _, sourceID := range source.ResourceIDs {
			switch source.Platform {
			case common.PlatformGithub:
				summary, err := source.GitHubMetrics(sourceID.ID, httpClient)
				if err != nil {
					cfg.logger.Errorf("fetching download metrics with id '%s' of platform '%s' errored with: '%v'", sourceID, source.Platform, err)
				}

				insight = append(insight, Insight{
					Platform: source.Platform,
					ID:       sourceID.GetID(),
					Summary:  summary,
				})
			case common.PlatformTerraform:
				summary, err := source.TerraformMetrics(sourceID.Name, sourceID.ID, httpClient)
				if err != nil {
					cfg.logger.Errorf("fetching download metrics with id '%s' of platform '%s' errored with: '%v'", sourceID, source.Platform, err)
				}

				insight = append(insight, Insight{
					Platform: source.Platform,
					ID:       sourceID.GetID(),
					Summary:  summary,
				})
			}
		}
	}

	return insight
}
