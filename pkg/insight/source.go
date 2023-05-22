package insight

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jinzhu/copier"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/common"
	"github.com/sirupsen/logrus"
)

type Source struct {
	Platform   string   `json:"platform,omitempty" yaml:"platform,omitempty"`
	BaseURL    string   `json:"base_url,omitempty" yaml:"base_url,omitempty"`
	IDs        []string `json:"ids,omitempty" yaml:"ids,omitempty"`
	CaFilePath string   `json:"ca_file_path,omitempty" yaml:"ca_file_path,omitempty"`
	Auth       Auth     `json:"auth,omitempty" yaml:"auth,omitempty"`
	logger     *logrus.Logger
}

// NewClient returns new instance of httpClient when invoked.
func (s *Source) NewClient(caContent []byte) *resty.Client {
	httpClient := resty.New()
	httpClient.SetLogger(s.logger)
	httpClient.SetRetryCount(common.DefaultRetryCount)
	httpClient.SetRetryWaitTime(common.DefaultRetryWaitTime * time.Second)
	if s.logger.Level == logrus.DebugLevel {
		httpClient.SetDebug(true)
	}

	if !reflect.DeepEqual(s.Auth, Auth{}) {
		s.Auth.SetAuth(httpClient)
	}

	// setting authorization
	httpClient.SetBaseURL(s.BaseURL)

	if len(caContent) != 0 {
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caContent)
		httpClient.SetTLSClientConfig(&tls.Config{RootCAs: certPool}) //nolint:gosec
	} else {
		httpClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	}

	return httpClient
}

func (s *Source) SetLogger(log *logrus.Logger) {
	s.logger = log
}

func (s *Source) Metrics(sourceID string, httpClient *resty.Client) (any, error) {
	newClient := resty.New()
	if err := copier.CopyWithOption(newClient, httpClient, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var endPoint string

	switch s.Platform {
	case common.PlatformGithub:
		endPoint = filepath.Join(sourceID, "releases")
	case common.PlatformTerraform:
		endPoint = filepath.Join(sourceID, "/downloads/summary")
	}

	response, err := newClient.R().Get(endPoint)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("fetching download metrics returned non ok status code '%d' with body: ", response.StatusCode())
	}

	return s.ReadResponse(response.Body())
}

func (s *Source) ReadResponse(body []byte) (any, error) {
	switch s.Platform {
	case common.PlatformGithub:
		var gitRelease []GitRelease
		if err := json.Unmarshal(body, &gitRelease); err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return gitRelease, nil
	case common.PlatformTerraform:
		var downloadSummary ProviderDownloadSummary
		if err := json.Unmarshal(body, &downloadSummary); err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		return downloadSummary, nil
	default:
		return nil, fmt.Errorf("unknown platform")
	}
}

func (s *Source) GetPlatform() (string, error) {
	if strings.Contains(s.BaseURL, common.PlatformTerraform) {
		return common.PlatformTerraform, nil
	}
	if strings.Contains(s.BaseURL, common.PlatformGithub) {
		return common.PlatformGithub, nil
	}

	return "", fmt.Errorf("unknown platform, either insighter doesnot support the platform currently or platform doesnot exists: %s", s.BaseURL)
}
