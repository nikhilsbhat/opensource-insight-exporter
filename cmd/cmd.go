package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/common"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/exporter"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/insight"
	"github.com/nikhilsbhat/opensource-insight-exporter/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
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
		Authors: []*cli.Author{
			{
				Name:  "Nikhil Bhat",
				Email: "nikhilsbhat93@gmail.com",
			},
		},
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
	logger := logrus.New()
	logger.SetLevel(gocd.GetLoglevel(context.String(flagLogLevel)))
	logger.WithField(common.OpenSourceInsightName, true)
	logger.SetFormatter(&logrus.JSONFormatter{})

	config := insight.Config{
		LogLevel: context.String(flagLogLevel),
	}

	cfg, err := insight.GetConfig(config, context.String(flagConfigPath))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	cfg.SetLogger(logger)

	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT) //nolint:govet

	exporterInsight := exporter.NewExporter(logger, *cfg)
	prometheus.MustRegister(exporterInsight)

	exporterPort := context.Int(flagExporterPort)
	appGraceDuration := context.Duration(flagGraceDuration)
	exporterEndpoint := context.String(flagExporterEndpoint)

	// listens to terminate signal
	go func() {
		sig := <-sigChan
		logger.Infof("caught signal %v: terminating in %v", sig, appGraceDuration)
		time.Sleep(appGraceDuration)
		logger.Infof("terminate opensource-insight-exporter running on port: %d", exporterPort)
		os.Exit(0)
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(getRedirectData(exporterEndpoint)))
	})

	logger.Infof("metrics will be exposed on port: %d", exporterPort)
	logger.Infof("metrics will be exposed on endpoint: %s", exporterEndpoint)

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
