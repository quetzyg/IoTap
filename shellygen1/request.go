package shellygen1

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	iotune "github.com/Stowify/IoTune"
	"github.com/Stowify/IoTune/device"
)

// buildURL for Shelly Gen1 requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// Create a Shelly Gen1 compliant request.
func request(dev device.Resource, path string, params any) (*http.Request, error) {
	values := structToValues(params)

	r, err := http.NewRequest(http.MethodPost, buildURL(dev.IP(), path), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	r.Header.Set(iotune.ContentTypeHeader, iotune.URLEncodedFormMimeType)

	return r, nil
}

// structToValues returns the url.Values representation of a struct. Unfortunately, the
// Shelly Gen1 API doesn't support JSON requests, only HTTP GET with a query-string in
// the URL or an HTTP POST with an application/x-www-form-urlencoded payload.
// Read more at: https://shelly-api-docs.shelly.cloud/gen1/#common-http-api
func structToValues(cfg any) url.Values {
	cfgVal := reflect.Indirect(reflect.ValueOf(cfg))
	cfgTyp := cfgVal.Type()

	var values = url.Values{}

	for i := 0; i < cfgVal.NumField(); i++ {
		fieldValue := cfgVal.Field(i)

		// Ignore nil pointers
		if cfgTyp.Field(i).Type.Kind() == reflect.Ptr && !fieldValue.Elem().CanAddr() {
			continue
		}

		value := fmt.Sprint(reflect.Indirect(fieldValue).Interface())

		// Convert the Schedule Rules array to string (i.e. CSV), since that's what
		// the Shelly API expects. Otherwise, we'll get "Bad schedule rules!" errors
		// when passing URL encoded arrays.
		key := strings.TrimSuffix(cfgTyp.Field(i).Tag.Get("json"), ",omitempty")

		if key == "schedule_rules" {
			value = strings.Trim(value, "[]")
			value = strings.ReplaceAll(value, " ", ",")
		}

		values.Add(key, value)
	}

	return values
}
