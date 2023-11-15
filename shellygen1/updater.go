package shellygen1

import (
	"net/http"
	"net/url"
	"strings"

	iotune "github.com/Stowify/IoTune"
)

// UpdateRequest returns a device firmware update HTTP request.
// See: https://shelly-api-docs.shelly.cloud/gen1/#ota
func (d *Device) UpdateRequest() (*http.Request, error) {
	values := url.Values{}
	values.Add("update", "true")

	r, err := http.NewRequest(http.MethodPost, buildURL(d.ip, updatePath), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	r.Header.Set(iotune.ContentTypeHeader, iotune.URLEncodedFormMimeType)

	return r, nil
}
