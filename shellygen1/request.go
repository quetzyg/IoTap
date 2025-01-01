package shellygen1

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/quetzyg/IoTap/httpclient"
)

// buildURL for Shelly Gen1 requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// Create a Shelly Gen1 compliant request.
func request(dev *Device, path string, params any) (*http.Request, error) {
	values, ok := params.(url.Values)
	if !ok && params != nil {
		values = settingsToValues(params)
	}

	if len(values) > 0 {
		path = fmt.Sprintf("%s?%s", path, values.Encode())
	}

	r, err := http.NewRequest(http.MethodGet, buildURL(dev.IP(), path), nil)
	if err != nil {
		return nil, err
	}

	r.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

	if dev.Secured() {
		return dev.SecureRequest(r)
	}

	return r, nil
}

// settingsToValues returns the url.Values representation of a *settings type. Unfortunately,
// the Shelly Gen1 API doesn't support JSON requests, only HTTP GET with a query-string URL
// or an HTTP POST with an application/x-www-form-urlencoded body payload.
// Read more at: https://shelly-api-docs.shelly.cloud/gen1/#common-http-api
func settingsToValues(params any) url.Values {
	m, ok := params.(*settings)
	if !ok {
		return nil
	}

	var values = url.Values{}

	for key, value := range *m {
		switch val := value.(type) {
		case []any:
			// Convert the Schedule Rules array to a CSV
			// since that's what the Shelly API expects.
			if key == "schedule_rules" {
				var rules []string
				for _, rule := range val {
					rules = append(rules, fmt.Sprint(rule))
				}
				values.Add(key, strings.Join(rules, ","))
				continue
			}

			// Handle other slices as usual.
			for _, v := range val {
				values.Add(key, fmt.Sprint(v))
			}

		case nil:
			values.Add(key, "null")

		default:
			values.Add(key, fmt.Sprint(val))
		}
	}

	return values
}
