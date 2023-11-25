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
			name: "success: single level",
			key:  "foo",
			value: map[string]any{
				"foo": 123,
			},
			exists: true,
		},
		{
			name: "failure: single level #1",
			key:  "baz",
			value: map[string]any{
				"foo": "value",
			},
		},
		{
			name: "failure: single level #2",
			key:  "foo.bar",
			value: map[string]any{
				"foo": true,
			},
		},
		{
			name: "success: multi level #1",
			key:  "foo",
			value: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{
							"qux": map[string]any{
								"quux": []int{},
							},
						},
					},
				},
			},
			exists: true,
		},
		{
			name: "success: multi level #2",
			key:  "foo.bar.baz",
			value: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{
							"qux": map[string]any{
								"quux": "value",
							},
						},
					},
				},
			},
			exists: true,
		},
		{
			name: "success: multi level #3",
			key:  "foo.bar.baz.qux.quux",
			value: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{
							"qux": map[string]any{
								"quux": false,
							},
						},
					},
				},
			},
			exists: true,
		},
		{
			name: "failure: multi level #1",
			key:  "xyz",
			value: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{
							"qux": map[string]any{
								"quux": 789,
							},
						},
					},
				},
			},
		},
		{
			name: "failure: multi level #2",
			key:  "foo.bar.xyz",
			value: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{
							"qux": map[string]any{
								"quux": 42,
							},
						},
					},
				},
			},
		},
		{
			name: "failure: multi level #3",
			key:  "foo.bar.baz.qux.quux.xyz",
			value: map[string]any{
				"foo": map[string]any{
					"bar": map[string]any{
						"baz": map[string]any{
							"qux": map[string]any{
								"quux": true,
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			exists := KeyExists(test.value, test.key)
			if exists != test.exists {
				t.Fatalf("expected %t, got %t", test.exists, exists)
			}
		})
	}
}
