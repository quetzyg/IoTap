package device

import (
	"encoding/json"
	"errors"
	"net"
	"testing"
)

func TestStrategy_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
		err  error
	}{
		{
			name: "failure: invalid JSON",
			err:  &json.SyntaxError{},
		},
		{
			name: "failure: undefined strategy mode #1",
			data: "{}",
			err:  errStrategyModeUndefined,
		},
		{
			name: "failure: undefined strategy mode #2",
			data: `{"mode":""}`,
			err:  errStrategyModeUndefined,
		},
		{
			name: "failure: invalid strategy mode",
			data: `{"mode":"foo"}`,
			err:  errStrategyModeInvalid,
		},
		{
			name: "success: whitelist strategy",
			data: `{"mode":"whitelist"}`,
		},
		{
			name: "success: blacklist strategy",
			data: `{"mode":"blacklist"}`,
		},
		{
			name: "failure: blacklist strategy invalid device",
			data: `{
			  "mode": "blacklist",
			  "devices": [
				"foo"
			  ]
			}`,
			err: &net.AddrError{},
		},
		{
			name: "success: whitelist strategy invalid device",
			data: `{
			  "mode": "whitelist",
			  "devices": [
				"14:06:12:DC:7A:F0"
			  ]
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(test.data), &Strategy{})

			switch v := test.err.(type) {
			case nil:
				if errors.Is(err, v) {
					return
				}

			case *json.SyntaxError:
				var se *json.SyntaxError
				if errors.As(err, &se) {
					return
				}
			case *net.AddrError:
				var ae *net.AddrError
				if errors.As(err, &ae) {
					return
				}

			default:
				if errors.Is(err, test.err) {
					return
				}
			}

			t.Fatalf("expected %v, got %v", test.err, err)
		})
	}
}
