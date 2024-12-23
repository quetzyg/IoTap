package device

import "testing"

func TestMacify(t *testing.T) {
	tests := []struct {
		name     string
		mac      string
		expected string
	}{
		{
			name:     "return invalid mac address format",
			mac:      "ab12cd",
			expected: "ab12cd",
		},
		{
			name:     "return valid mac address format",
			mac:      "ab12cd34ef09",
			expected: "ab:12:cd:34:ef:09",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Macify(test.mac)
			if got != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, got)
			}
		})
	}
}
