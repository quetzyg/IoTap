package device

import (
	"reflect"
	"testing"
)

func TestRegisterProber(t *testing.T) {
	if len(proberRegistry) != 0 {
		t.Fatal("Prober registry should be empty")
	}

	RegisterProber("foo", func() Prober {
		return &prober{}
	})

	if len(proberRegistry) != 1 {
		t.Fatalf("Prober registry should have one registered prober, %d found", len(proberRegistry))
	}
}

func TestGetProbers(t *testing.T) {
	tests := []struct {
		name    string
		driver  string
		probers []Prober
	}{
		{
			name:   "return all probers",
			driver: AllDrivers,
			probers: []Prober{
				&prober{},
				&prober{},
				&prober{},
			},
		},
		{
			name:   "return single prober",
			driver: "bar",
			probers: []Prober{
				&prober{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proberRegistry = map[string]ProberProvider{
				"foo": func() Prober {
					return &prober{}
				},
				"bar": func() Prober {
					return &prober{}
				},
				"baz": func() Prober {
					return &prober{}
				},
			}

			probers := GetProbers(test.driver)

			if !reflect.DeepEqual(probers, test.probers) {
				t.Fatalf("expected %#v, got %#v", test.probers, probers)
			}
		})
	}
}
