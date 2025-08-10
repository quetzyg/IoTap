package shellygen2

import "testing"

func TestConfig_Driver(t *testing.T) {
	driver := (&Config{}).Driver()

	if driver != Driver {
		t.Fatalf("expected %q, got %q", Driver, driver)
	}
}

func TestConfig_Empty(t *testing.T) {
	tests := []struct {
		cfg   *Config
		name  string
		empty bool
	}{
		{
			name:  "empty config",
			cfg:   &Config{},
			empty: true,
		},
		{
			name: "not empty config",
			cfg: &Config{
				Sys: &settings{},
			},
			empty: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.cfg.Empty() != test.empty {
				t.Fatalf("expected empty to be %t, got %t", test.empty, test.cfg.Empty())
			}
		})
	}
}
