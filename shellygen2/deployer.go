package shellygen2

import (
	"net/http"

	"github.com/quetzyg/IoTap/device"
	"github.com/quetzyg/IoTap/httpclient"
)

const chunkSize = 1024

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

// fetchScripts constructs and returns a slice of *script resources associated with the device.
func (d *Device) fetchScripts(client *http.Client) ([]*script, error) {
	// List all the IoT device scripts
	// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptlist
	r, err := request(d, "Script.List", nil)
	if err != nil {
		return nil, err
	}

	resp := &listResponse{}

	if err = httpclient.Dispatch(client, r, d, resp); err != nil {
		return nil, err
	}

	return resp.Result.Scripts, nil
}

// DeployRequests creates an ordered slice of *http.Request objects for deploying device scripts.
func (d *Device) DeployRequests(client *http.Client, dep *device.Deployment) ([]*http.Request, error) {
	// Check if a deployment policy is set and enforce it
	if dep.Policy != nil && dep.Policy.IsExcluded(d) {
		return nil, device.ErrPolicyExcluded
	}

	var requests []*http.Request

	// Delete any existing scripts
	// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptdelete
	scripts, err := d.fetchScripts(client)
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

	// Deploy scripts
	for id, s := range dep.Scripts {
		// Create script
		// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptcreate
		r, err := request(d, "Script.Create", map[string]any{"name": s.Name()})
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)

		// Upload code in chunks
		// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptputcode
		for start := 0; start < s.Length(); start += chunkSize {
			end := start + chunkSize
			if end > s.Length() {
				end = s.Length()
			}

			r, err = request(d, "Script.PutCode", map[string]any{
				"id":     id + 1,
				"append": start != 0,
				"code":   string(s.Code()[start:end]),
			})
			if err != nil {
				return nil, err
			}
			requests = append(requests, r)
		}

		// Enable script
		// See: https://shelly-api-docs.shelly.cloud/gen2/ComponentsAndServices/Script#scriptsetconfig
		r, err = request(d, "Script.SetConfig", map[string]any{
			"id": id + 1,
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
		r, err = request(d, "Script.Start", map[string]any{"id": id + 1})
		if err != nil {
			return nil, err
		}
		requests = append(requests, r)
	}

	// Reboot request
	r, err := request(d, "Shelly.Reboot", nil)
	if err != nil {
		return nil, err
	}

	return append(requests, r), nil
}
