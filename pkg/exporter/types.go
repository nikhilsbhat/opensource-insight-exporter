package exporter

import (
	"sync"

	"github.com/go-kit/log"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/common"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/insight"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	config              insight.Config
	mutex               sync.Mutex
	logger              log.Logger
	downloadMetricCount *prometheus.GaugeVec
}

func NewExporter(logger log.Logger, config insight.Config) *Exporter {
	return &Exporter{
		config: config,
		logger: logger,
		downloadMetricCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: common.OpenSourceInsightExporterName,
			Name:      common.DownloadMetrics,
			Help:      "downloads count of artifacts hosted  in multiple platforms",
		}, []string{"platform", "id", "name", "tag", "attributes"},
		),
	}
}
