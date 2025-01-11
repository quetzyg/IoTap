package device

import "testing"

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
