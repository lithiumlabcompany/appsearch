package appsearch

import (
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// Open APIClient with endpoint and key
// First parameter may be specified as URL with API key in authentication like:
// https://private-...@abcd.ent-search.eu-central-1.aws.cloud.es.io
// Second parameter is always interpreted as API key if specified
func Open(endpointAndKey ...string) (APIClient, error) {
	hostURL, token, authType, err := getHostURL(endpointAndKey)

	return &client{
		resty.New().
			SetHostURL(hostURL).
			SetAuthToken(token).
			SetAuthScheme(authType),
	}, err
}

func getHostURL(params []string) (hostURL string, token string, authType string, err error) {
	switch len(params) {
	case 2:
		hostURL, _, _, err = resolve(params[0])
		authType = "Bearer"
		token = params[1]
		return
	case 1:
		return resolve(params[0])
	default:
		return "", "", "", ErrInvalidParams
	}
}

func resolve(rawURL string) (hostURL string, token string, authType string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}
	if u.User != nil {
		username := u.User.Username()
		password, isBasicAuth := u.User.Password()
		if isBasicAuth {
			authType = "Basic"
			token = base64.StdEncoding.EncodeToString([]byte(
				fmt.Sprintf("%s:%s", u.User.Username(), password),
			))
		} else {
			authType = "Bearer"
			token = username
		}
		u.User = nil
	}
	hostURL = u.ResolveReference(&url.URL{Path: "/api/as/v1/"}).String()
	return
}
