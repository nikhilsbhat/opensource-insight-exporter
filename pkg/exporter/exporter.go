package exporter

import (
	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/common"
	insight2 "github.com/nikhilsbhat/opensource-insight-exporter/pkg/insight"
	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) collect(channel chan<- prometheus.Metric) {
	insights := e.config.GetInsight()

	for _, insight := range insights {
		switch insight.Platform {
		case common.PlatformGithub:
			summary, ok := insight.Summary.([]insight2.GitRelease)
			if !ok {
				level.Error(e.logger).Log(common.LogCategoryErr, "failed to typecast to []insight.GitRelease") //nolint:errcheck
			}

			for _, smry := range summary {
				for _, assetDownloadSummary := range smry.Assets {
					e.downloadMetricCount.WithLabelValues(
						insight.Platform,
						insight.ID,
						assetDownloadSummary.Name,
						smry.TagName, "").Set(common.Float(assetDownloadSummary.DownloadCount))
				}
			}
		case common.PlatformTerraform:
			summary, ok := insight.Summary.(insight2.ProviderDownloadSummary)
			if !ok {
				level.Error(e.logger).Log(common.LogCategoryErr, "failed to typecast to insight.ProviderDownloadSummary") //nolint:errcheck
			}

			attributes := summary.Data.Attributes
			e.downloadMetricCount.WithLabelValues(insight.Platform, insight.ID, "", "", "week").Set(common.Float(attributes.Week))
			e.downloadMetricCount.WithLabelValues(insight.Platform, insight.ID, "", "", "month").Set(common.Float(attributes.Month))
			e.downloadMetricCount.WithLabelValues(insight.Platform, insight.ID, "", "", "year").Set(common.Float(attributes.Year))
			e.downloadMetricCount.WithLabelValues(insight.Platform, insight.ID, "", "", "total").Set(common.Float(attributes.Total))
		}
	}

	e.downloadMetricCount.Collect(channel)
}

func (e *Exporter) Describe(channel chan<- *prometheus.Desc) {
	e.downloadMetricCount.Describe(channel)
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	e.collect(ch)
}
