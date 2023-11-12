package maputil

import "testing"

func TestKeyExists(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		value  map[string]any
		exists bool
	}{
		{
			name: "success: one level",
			key:  "foo",
			value: map[string]any{
				"foo": "bar",
			},
			exists: true,
		},
		{
			name: "failure: one level",
			key:  "baz",
			value: map[string]any{
				"foo": "bar",
			},
		},
		{
			name: "success: two levels",
			key:  "foo.bar",
			value: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
			exists: true,
		},
		{
			name: "failure: two levels",
			key:  "baz",
			value: map[string]any{
				"foo": map[string]any{
					"bar": "baz",
				},
			},
		},
		{
			name: "success: three levels",
			key:  "foo.bar.baz",
			value: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": "qux",
					},
				},
			},
			exists: true,
		},
		{
			name: "failure: three levels",
			key:  "baz.bar.foo",
			value: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": "qux",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exists := KeyExists(test.value, test.key)
			if exists != test.exists {
				t.Fatalf("expected %t when validating `%#v`, got %t", test.exists, test.value, exists)
			}
		})
	}
}
