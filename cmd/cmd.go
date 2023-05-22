package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-kit/log/level"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/common"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/exporter"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/insight"
	"github.com/nikhilsbhat/opensource-insight-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/urfave/cli/v2"
)

const (
	flagLogLevel         = "log-level"
	flagExporterPort     = "port"
	flagExporterEndpoint = "endpoint"
	flagCaPath           = "ca-path"
	flagConfigPath       = "config-file"
	flagGraceDuration    = "grace-duration"
)

const (
	defaultAppPort       = 8090
	defaultTimeout       = 30
	defaultGraceDuration = 5
)

var (
	redirectData = `<html>
			 <head><title>OpenSource Insight Exporter</title></head>
			 <body>
			 <h1>OpenSource Insight Exporter</h1>
			 <p><a href='%s'>Metrics</a></p>
			 </body>
			 </html>`
	sigChan = make(chan os.Signal)
)

// App returns the cli for opensource-insight-exporter.
func App() *cli.App {
	return &cli.App{
		Name:                 "opensource-insight-exporter",
		Usage:                "Utility to collect download metrics of open source projects hosted on multiple platforms",
		UsageText:            "opensource-insighter-exporter [flags]",
		EnableBashCompletion: true,
		HideHelp:             false,
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "version of the opensource-insight-exporter",
				Action:  version.AppVersion,
			},
		},
		Flags:  registerFlags(),
		Action: insightExporter,
	}
}

func registerFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    flagLogLevel,
			Usage:   "set log level for the opensource insight exporter",
			Aliases: []string{"log"},
			Value:   "info",
		},
		&cli.IntFlag{
			Name:    flagExporterPort,
			Usage:   "port on which the metrics to be exposed",
			Value:   defaultAppPort,
			Aliases: []string{"p"},
		},
		&cli.StringFlag{
			Name:    flagExporterEndpoint,
			Usage:   "path under which the metrics to be exposed",
			Value:   "/metrics",
			Aliases: []string{"e"},
		},
		&cli.StringFlag{
			Name:    flagCaPath,
			Usage:   "path to file containing CA information to make secure connections to various sources",
			Aliases: []string{"ca"},
		},
		&cli.StringFlag{
			Name:    flagConfigPath,
			Usage:   "path to file containing configurations for exporter",
			Aliases: []string{"c"},
			Value:   filepath.Join(os.Getenv("HOME"), fmt.Sprintf("%s.%s", common.OpenSourceInsightName, common.ExporterConfigFileExt)),
		},
		&cli.DurationFlag{
			Name:    flagGraceDuration,
			Usage:   "time duration to wait before stopping the service",
			Aliases: []string{"d"},
			Value:   time.Second * defaultGraceDuration,
		},
	}
}

func insightExporter(context *cli.Context) error {
	config := insight.Config{
		LogLevel: context.String(flagLogLevel),
	}

	cfg, err := insight.GetConfig(config, context.String(flagConfigPath))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	fmt.Println("final config", cfg)

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT) //nolint:govet

	promLogConfig := &promlog.Config{Level: &promlog.AllowedLevel{}, Format: &promlog.AllowedFormat{}}
	if err = promLogConfig.Level.Set(config.LogLevel); err != nil {
		return fmt.Errorf("configuring logger errored with: %w", err)
	}
	logger := promlog.New(promLogConfig)

	exporterInsight := exporter.NewExporter(logger, *cfg)
	prometheus.MustRegister(exporterInsight)

	exporterPort := context.Int(flagExporterPort)
	appGraceDuration := context.Duration(flagGraceDuration)
	exporterEndpoint := context.String(flagExporterEndpoint)

	// listens to terminate signal
	go func() {
		sig := <-sigChan
		level.Info(logger).Log("msg", fmt.Sprintf("caught signal %v: terminating in %v", sig, appGraceDuration)) //nolint:errcheck
		time.Sleep(appGraceDuration)
		level.Info(logger).Log("msg", fmt.Sprintf("terminate opensource-insight-exporter running on port: %d", exporterPort)) //nolint:errcheck
		os.Exit(0)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(getRedirectData(exporterEndpoint)))
	})

	level.Info(logger).Log(common.LogCategoryMsg, fmt.Sprintf("metrics will be exposed on port: %d", exporterPort))         //nolint:errcheck
	level.Info(logger).Log(common.LogCategoryMsg, fmt.Sprintf("metrics will be exposed on endpoint: %s", exporterEndpoint)) //nolint:errcheck

	http.Handle(exporterEndpoint, promhttp.Handler())

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", exporterPort),
		ReadHeaderTimeout: defaultTimeout * time.Second,
	}

	if err = server.ListenAndServe(); err != nil {
		return fmt.Errorf("starting server on specified port failed with: %w", err)
	}

	return nil
}

func getRedirectData(endpoint string) string {
	return fmt.Sprintf(redirectData, endpoint)
}
