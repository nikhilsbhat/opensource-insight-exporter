package insight

import (
	"log"

	"github.com/go-resty/resty/v2"
)

// Auth holds the authorization information of selected platforms.
type Auth struct {
	UserName    string `json:"username,omitempty" yaml:"username,omitempty"`
	Password    string `json:"password,omitempty" yaml:"password,omitempty"`
	BearerToken string `json:"bearer_token,omitempty" yaml:"bearer_token,omitempty"`
}

// SetAuth sets authorization for the http client if passed.
func (auth *Auth) SetAuth(httpClient *resty.Client) {
	if len(auth.BearerToken) != 0 {
		httpClient.SetAuthToken(auth.BearerToken)

		return
	}

	if len(auth.UserName) != 0 && len(auth.Password) != 0 {
		httpClient.SetBasicAuth(auth.UserName, auth.Password)

		return
	}

	log.Println("no auth specified, API calls would be made with no authentication")
}
