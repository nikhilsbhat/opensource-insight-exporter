package common

const (
	OpenSourceInsightName         = "opensource-insight-exporter"
	ExporterConfigFileExt         = "yaml"
	OpenSourceInsightExporterName = "opensource_insight_exporter"
	DownloadMetrics               = "downloads_count"
	LogCategoryMsg                = "msg"
	LogCategoryErr                = "err"
)

const (
	DefaultRetryCount    = 5
	DefaultRetryWaitTime = 5
)

const (
	PlatformTerraform = "terraform"
	PlatformGithub    = "github"
)

func Float(value interface{}) float64 {
	switch value.(type) {
	case int64:
		return value.(float64) //nolint:forcetypeassert
	case string:
		return float64(0)
	default:
		return float64(value.(int)) //nolint:forcetypeassert
	}
}
