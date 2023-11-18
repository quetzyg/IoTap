package shellygen2

import (
	"net/http"

	iotune "github.com/Stowify/IoTune"
	"github.com/Stowify/IoTune/device"
)

const (
	defaultID = 1
	chunkSize = 1024
)

// Basic script resource representation.
type script struct {
	ID int `json:"id"`
}

// listResponse holds the result of a Script.List method request.
type listResponse struct {
	Result struct {
		Scripts []*script `json:"scripts"`
	} `json:"result"`
}

// scripts returns all script resources of the device.
func (d *Device) scripts() ([]*script, error) {
	// List all the IoT device scripts
	// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptlist
	r, err := request(d, "Script.List", nil)
	if err != nil {
		return nil, err
	}

	resp := &listResponse{}
	err = iotune.Dispatch(&http.Client{}, r, resp)
	if err != nil {
		return nil, err
	}

	return resp.Result.Scripts, nil
}

// ScriptRequests generates a slice of *http.Requests that are to be executed in order to set an IoT device script.
func (d *Device) ScriptRequests(script *device.IoTScript) ([]*http.Request, error) {
	var requests []*http.Request

	// Delete any existing scripts
	// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptdelete
	scripts, err := d.scripts()
	if err != nil {
		return nil, err
	}

	for _, s := range scripts {
		r, err := request(d, "Script.Delete", map[string]any{"id": s.ID})
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	// Create script
	// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptcreate
	r, err := request(d, "Script.Create", map[string]any{"name": script.Name()})
	if err != nil {
		return nil, err
	}
	requests = append(requests, r)

	// Upload code in chunks
	// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptputcode
	for start := 0; start < script.Length(); start += chunkSize {
		end := start + chunkSize
		if end > script.Length() {
			end = script.Length()
		}

		r, err = request(d, "Script.PutCode", map[string]any{
			"id":     defaultID,
			"append": start != 0,
			"code":   string(script.Code()[start:end]),
		})
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	// Enable script
	// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptsetconfig
	r, err = request(d, "Script.SetConfig", map[string]any{
		"id": defaultID,
		"config": map[string]any{
			"enable": true,
		},
	})
	if err != nil {
		return nil, err
	}
	requests = append(requests, r)

	// Start script
	// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptstart
	r, err = request(d, "Script.Start", map[string]any{"id": defaultID})
	if err != nil {
		return nil, err
	}

	return append(requests, r), nil
}
