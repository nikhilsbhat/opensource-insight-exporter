// Package version powers the versioning of opensource-insight-exporter.
package version

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	// Version specifies the version of the application and cannot be changed by end user.
	Version string

	// Env tells end user that what variant (here we use the name of the git branch to make it simple)
	// of application is he using.
	Env string

	// BuildDate of the app.
	BuildDate string
	// GoVersion represents golang version used.
	GoVersion string
	// Platform is the combination of OS and Architecture for which the binary is built for.
	Platform string
	// Revision represents the git revision used to build the current version of app.
	Revision string
)

// BuildInfo represents version of utility.
type BuildInfo struct {
	Version     string `json:"version,omitempty" yaml:"version,omitempty"`
	Revision    string `json:"revision,omitempty" yaml:"revision,omitempty"`
	Environment string `json:"environment,omitempty" yaml:"environment,omitempty"`
	BuildDate   string `json:"buildDate,omitempty" yaml:"buildDate,omitempty"`
	GoVersion   string `json:"goVersion,omitempty" yaml:"goVersion,omitempty"`
	Platform    string `json:"platform,omitempty" yaml:"platform,omitempty"`
}

// GetBuildInfo return the version and other build info of the application.
func GetBuildInfo() BuildInfo {
	if strings.ToLower(Env) != "production" {
		Env = "alfa"
	}

	return BuildInfo{
		Version:     Version,
		Revision:    Revision,
		Environment: Env,
		Platform:    Platform,
		BuildDate:   BuildDate,
		GoVersion:   GoVersion,
	}
}

func AppVersion(_ *cli.Context) error {
	buildInfo, err := json.Marshal(GetBuildInfo())
	if err != nil {
		return fmt.Errorf("fetching app version errored with: %w", err)
	}
	fmt.Println(string(buildInfo))

	return nil
}
