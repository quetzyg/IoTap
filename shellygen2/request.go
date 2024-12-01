package shellygen2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Stowify/IoTap/httpclient"
)

const rpcPath = "rpc"

// rpcRequest represents a command for an IoT device to execute.
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
func request(dev *Device, method string, params any) (*http.Request, error) {
	rpc := &rpcRequest{
		Source: "IoTap",
		Method: method,
	}

	if params != nil {
		rpc.Parameters = params
	}

	body, err := json.Marshal(rpc)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, buildURL(dev.IP(), rpcPath), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	r.Header.Set(httpclient.ContentTypeHeader, httpclient.JSONMimeType)

	return r, nil
}
