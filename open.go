package appsearch

import (
	"net/url"

	"github.com/go-resty/resty/v2"
)

// Open APIClient with endpoint and key
// First parameter may be specified as URL with API key in authentication like:
// https://private-...@abcd.ent-search.eu-central-1.aws.cloud.es.io
// Second parameter is always interpreted as API key if specified
func Open(endpointAndKey ...string) (*client, error) {
	hostURL, key, err := getHostURL(endpointAndKey)

	return &client{
		resty.New().
			SetHostURL(hostURL).
			SetAuthToken(key).
			SetAuthScheme("Bearer"),
	}, err
}

func getHostURL(params []string) (hostUrl string, key string, err error) {
	switch len(params) {
	case 2:
		hostUrl, _, err = resolve(params[0])
		key = params[1]
		return hostUrl, key, err
	case 1:
		return resolve(params[0])
	default:
		return "", "", ErrInvalidParams
	}
}

func resolve(rawUrl string) (hostUrl string, key string, err error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return
	}
	if u.User != nil {
		key = u.User.String()
		u.User = nil
	}
	hostUrl = u.ResolveReference(&url.URL{Path: "/api/as/v1/"}).String()
	return
}
