package insight

import (
	"os"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/common"
	"github.com/sirupsen/logrus"
)

type Insight struct {
	Platform string
	ID       string
	Summary  any
}

func (c *Config) GetInsight() []Insight {
	logger := logrus.New()
	logger.SetLevel(gocd.GetLoglevel(c.LogLevel))
	logger.WithField(common.OpenSourceInsightName, true)
	logger.SetFormatter(&logrus.JSONFormatter{})

	insight := make([]Insight, 0)
	for _, source := range c.Sources {
		source.SetLogger(logger)

		caContent := make([]byte, 0)
		if len(source.CaFilePath) != 0 {
			ca, err := os.ReadFile(source.CaFilePath)
			if err != nil {
				logger.Errorf("reading ca file errored with %v", err)
			}
			caContent = ca
		}

		httpClient := source.NewClient(caContent)

		for _, sourceID := range source.IDs {
			summary, err := source.Metrics(sourceID, httpClient)
			if err != nil {
				logger.Errorf("fetching download metrics with id '%s' of platform '%s' errored with: '%v'", sourceID, source.Platform, err)
			}
			insight = append(insight, Insight{
				Platform: source.Platform,
				ID:       sourceID,
				Summary:  summary,
			})
		}
	}

	return insight
}
