package shellygen2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	iotune "github.com/Stowify/IoTune"
	"github.com/Stowify/IoTune/device"
)

const rpcPath = "rpc"

// rpcRequest for the RPC endpoint.
// See: https://shelly-api-docs.shelly.cloud/gen2/General/RPCProtocol
type rpcRequest struct {
	ID         int    `json:"id"`
	Source     string `json:"src"`
	Method     string `json:"method"`
	Parameters any    `json:"params,omitempty"`
}

// buildURL for Shelly Gen2 requests.
func buildURL(ip net.IP, path string) string {
	return fmt.Sprintf("http://%s/%s", ip.String(), strings.TrimPrefix(path, "/"))
}

// Create a Shelly Gen2 compliant request.
func request(dev device.Resource, method string, params any) (*http.Request, error) {
	req := &rpcRequest{
		Source: "IoTune",
		Method: method,
	}

	if params != nil {
		req.Parameters = params
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, buildURL(dev.IP(), rpcPath), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	r.Header.Set(iotune.ContentTypeHeader, iotune.JSONMimeType)

	return r, nil
}
