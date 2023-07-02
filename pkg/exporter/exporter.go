package exporter

import (
	"time"

	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/common"
	insight2 "github.com/nikhilsbhat/opensource-insight-exporter/pkg/insight"
	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) collect(channel chan<- prometheus.Metric) {
	timeBeforeCalculation := time.Now()
	insights := e.config.GetInsight()

	timeAfterCalculation := time.Now()
	timeConsumed := timeAfterCalculation.Sub(timeBeforeCalculation)
	e.logger.Infof("time taken to collect the metrics is: %v", timeConsumed)

	for _, insight := range insights {
		switch insight.Platform {
		case common.PlatformGithub:
			summary, ok := insight.Summary.([]insight2.GitRelease)
			if !ok {
				e.logger.Error("failed to typecast to []insight.GitRelease")
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
			summary, ok := insight.Summary.([]insight2.ProviderDownloadSummary)
			if !ok {
				e.logger.Errorf("failed to typecast to insight.ProviderDownloadSummary")
			}

			for _, smry := range summary {
				e.downloadMetricCount.WithLabelValues(insight.Platform, insight.ID, smry.Data.ID, "", "week").Set(common.Float(smry.Data.Attributes.Week))
				e.downloadMetricCount.WithLabelValues(insight.Platform, insight.ID, smry.Data.ID, "", "month").Set(common.Float(smry.Data.Attributes.Month))
				e.downloadMetricCount.WithLabelValues(insight.Platform, insight.ID, smry.Data.ID, "", "year").Set(common.Float(smry.Data.Attributes.Year))
				e.downloadMetricCount.WithLabelValues(insight.Platform, insight.ID, smry.Data.ID, "", "total").Set(common.Float(smry.Data.Attributes.Total))
			}
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
