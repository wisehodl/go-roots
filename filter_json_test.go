package roots

import (
	"encoding/json"
	//"fmt"
	"reflect"
	"testing"
)

// Types

type FilterMarshalTestCase struct {
	name     string
	filter   Filter
	expected string
}

type FilterUnmarshalTestCase struct {
	name     string
	input    string
	expected Filter
}

type FilterRoundTripTestCase struct {
	name   string
	filter Filter
}

// Test Cases

var marshalTestCases = []FilterMarshalTestCase{
	{
		name:     "empty filter",
		filter:   Filter{},
		expected: `{}`,
	},

	// ID cases
	{
		name:     "nil IDs",
		filter:   Filter{IDs: nil},
		expected: `{}`,
	},

	{
		name:     "empty IDs",
		filter:   Filter{IDs: []string{}},
		expected: `{"ids":[]}`,
	},

	{
		name:     "populated IDs",
		filter:   Filter{IDs: []string{"abc", "123"}},
		expected: `{"ids":["abc","123"]}`,
	},

	// Author cases
	{
		name:     "nil Authors",
		filter:   Filter{Authors: nil},
		expected: `{}`,
	},

	{
		name:     "empty Authors",
		filter:   Filter{Authors: []string{}},
		expected: `{"authors":[]}`,
	},

	{
		name:     "populated Authors",
		filter:   Filter{Authors: []string{"abc", "123"}},
		expected: `{"authors":["abc","123"]}`,
	},

	// Kind cases
	{
		name:     "nil Kinds",
		filter:   Filter{Kinds: nil},
		expected: `{}`,
	},

	{
		name:     "empty Kinds",
		filter:   Filter{Kinds: []int{}},
		expected: `{"kinds":[]}`,
	},

	{
		name:     "populated Kinds",
		filter:   Filter{Kinds: []int{1, 20001}},
		expected: `{"kinds":[1,20001]}`,
	},

	// Since cases
	{
		name:     "nil Since",
		filter:   Filter{Since: nil},
		expected: `{}`,
	},

	{
		name:     "populated Since",
		filter:   Filter{Since: intPtr(1000)},
		expected: `{"since":1000}`,
	},

	// Until cases
	{
		name:     "nil Until",
		filter:   Filter{Until: nil},
		expected: `{}`,
	},

	{
		name:     "populated Until",
		filter:   Filter{Until: intPtr(1000)},
		expected: `{"until":1000}`,
	},

	// Limit cases
	{
		name:     "nil Limit",
		filter:   Filter{Limit: nil},
		expected: `{}`,
	},

	{
		name:     "populated Limit",
		filter:   Filter{Limit: intPtr(100)},
		expected: `{"limit":100}`,
	},

	// All standard fields
	{
		name: "all standard fields",
		filter: Filter{
			IDs:     []string{"abc", "123"},
			Authors: []string{"def", "456"},
			Kinds:   []int{1, 200, 3000},
			Since:   intPtr(1000),
			Until:   intPtr(2000),
			Limit:   intPtr(100),
		},
		expected: `{"ids":["abc","123"],"authors":["def","456"],"kinds":[1,200,3000],"since":1000,"until":2000,"limit":100}`,
	},

	{
		name:     "mixed fields",
		filter:   Filter{IDs: nil, Authors: []string{}, Kinds: []int{1}},
		expected: `{"authors":[],"kinds":[1]}`,
	},

	// Tags
	{
		name:     "nil tags map",
		filter:   Filter{Tags: nil},
		expected: `{}`,
	},

	{
		name: "single-letter tag",
		filter: Filter{Tags: map[string][]string{
			"e": {"event1"},
		}},
		expected: `{"#e":["event1"]}`,
	},

	{
		name: "multi-letter tag",
		filter: Filter{Tags: map[string][]string{
			"emoji": {"🔥", "💧"},
		}},
		expected: `{"#emoji":["🔥","💧"]}`,
	},

	{
		name: "empty tag array",
		filter: Filter{Tags: map[string][]string{
			"p": {},
		}},
		expected: `{"#p":[]}`,
	},

	{
		name: "multiple tags",
		filter: Filter{Tags: map[string][]string{
			"e": {"event1", "event2"},
			"p": {"pubkey1", "pubkey2"},
		}},
		expected: `{"#e":["event1","event2"],"#p":["pubkey1","pubkey2"]}`,
	},

	// Extensions
	{
		name: "simple extension",
		filter: Filter{
			Extensions: map[string]json.RawMessage{
				"search": json.RawMessage(`"query"`),
			},
		},
		expected: `{"search":"query"}`,
	},

	{
		name: "extension with nested object",
		filter: Filter{
			Extensions: map[string]json.RawMessage{
				"meta": json.RawMessage(`{"author":"alice","score":99}`),
			},
		},
		expected: `{"meta":{"author":"alice","score":99}}`,
	},

	{
		name: "extension with nested array",
		filter: Filter{
			Extensions: map[string]json.RawMessage{
				"items": json.RawMessage(`[1,2,3]`),
			},
		},
		expected: `{"items":[1,2,3]}`,
	},

	{
		name: "extension with complex nested structure",
		filter: Filter{
			Extensions: map[string]json.RawMessage{
				"data": json.RawMessage(`{"users":[{"id":1}],"count":5}`),
			},
		},
		expected: `{"data":{"users":[{"id":1}],"count":5}}`,
	},

	{
		name: "multiple extensions",
		filter: Filter{
			Extensions: map[string]json.RawMessage{
				"search": json.RawMessage(`"x"`),
				"depth":  json.RawMessage(`3`),
			},
		},
		expected: `{"search":"x","depth":3}`,
	},

	// Extension Collisions
	{
		name: "extension collides with standard field - IDs",
		filter: Filter{
			IDs: []string{"real"},
			Extensions: map[string]json.RawMessage{
				"ids": json.RawMessage(`["fake"]`),
			},
		},
		expected: `{"ids":["real"]}`,
	},

	{
		name: "extension collides with standard field - Since",
		filter: Filter{
			Since: intPtr(100),
			Extensions: map[string]json.RawMessage{
				"since": json.RawMessage(`999`),
			},
		},
		expected: `{"since":100}`,
	},

	{
		name: "extension collides with multiple standard fields",
		filter: Filter{
			Authors: []string{"a"},
			Kinds:   []int{1},
			Extensions: map[string]json.RawMessage{
				"authors": json.RawMessage(`["b"]`),
				"kinds":   json.RawMessage(`[2]`),
			},
		},
		expected: `{"authors":["a"],"kinds":[1]}`,
	},

	{
		name: "extension collides with tag field - #e",
		filter: Filter{
			Extensions: map[string]json.RawMessage{
				"#e": json.RawMessage(`["fakeevent"]`),
			},
		},
		expected: `{}`,
	},

	{
		name: "extension collides with standard and tag fields",
		filter: Filter{
			Authors: []string{"realauthor"},
			Tags: map[string][]string{
				"e": {"realevent"},
			},
			Extensions: map[string]json.RawMessage{
				"authors": json.RawMessage(`["fakeauthor"]`),
				"#e":      json.RawMessage(`["fakeevent"]`),
			},
		},
		expected: `{"authors":["realauthor"],"#e":["realevent"]}`,
	},

	// Kitchen Sink
	{
		name: "filter with all field types",
		filter: Filter{
			IDs:   []string{"x"},
			Since: intPtr(100),
			Tags: map[string][]string{
				"e": {"y"},
			},
			Extensions: map[string]json.RawMessage{
				"search": json.RawMessage(`"z"`),
				"ids":    json.RawMessage(`["fakeid"]`),
			},
		},
		expected: `{"ids":["x"],"since":100,"#e":["y"],"search":"z"}`,
	},
}

var unmarshalTestCases = []FilterUnmarshalTestCase{
	{
		name:     "empty object",
		input:    `{}`,
		expected: Filter{},
	},

	// ID cases
	{
		name:     "null IDs",
		input:    `{"ids": null}`,
		expected: Filter{IDs: nil},
	},

	{
		name:     "empty IDs",
		input:    `{"ids": []}`,
		expected: Filter{IDs: []string{}},
	},

	{
		name:     "populated IDs",
		input:    `{"ids": ["abc","123"]}`,
		expected: Filter{IDs: []string{"abc", "123"}},
	},

	// Author cases
	{
		name:     "null Authors",
		input:    `{"authors": null}`,
		expected: Filter{Authors: nil},
	},

	{
		name:     "empty Authors",
		input:    `{"authors": []}`,
		expected: Filter{Authors: []string{}},
	},

	{
		name:     "populated Authors",
		input:    `{"authors": ["abc","123"]}`,
		expected: Filter{Authors: []string{"abc", "123"}},
	},

	// Kind cases
	{
		name:     "null Kinds",
		input:    `{"kinds": null}`,
		expected: Filter{Kinds: nil},
	},

	{
		name:     "empty Kinds",
		input:    `{"kinds": []}`,
		expected: Filter{Kinds: []int{}},
	},

	{
		name:     "populated Kinds",
		input:    `{"kinds": [1,2,3]}`,
		expected: Filter{Kinds: []int{1, 2, 3}},
	},

	// Since cases
	{
		name:     "null Since",
		input:    `{"since": null}`,
		expected: Filter{Since: nil},
	},

	{
		name:     "populated Since",
		input:    `{"since": 1000}`,
		expected: Filter{Since: intPtr(1000)},
	},

	// Until cases
	{
		name:     "null Until",
		input:    `{"until": null}`,
		expected: Filter{Until: nil},
	},

	{
		name:     "populated Until",
		input:    `{"until": 1000}`,
		expected: Filter{Until: intPtr(1000)},
	},

	// Limit cases
	{
		name:     "null Limit",
		input:    `{"limit": null}`,
		expected: Filter{Limit: nil},
	},

	{
		name:     "populated Limit",
		input:    `{"limit": 1000}`,
		expected: Filter{Limit: intPtr(1000)},
	},

	// All standard fields
	{
		name:  "all standard fields",
		input: `{"ids":["abc","123"],"authors":["def","456"],"kinds":[1,200,3000],"since":1000,"until":2000,"limit":100}`,
		expected: Filter{
			IDs:     []string{"abc", "123"},
			Authors: []string{"def", "456"},
			Kinds:   []int{1, 200, 3000},
			Since:   intPtr(1000),
			Until:   intPtr(2000),
			Limit:   intPtr(100),
		},
	},

	{
		name:     "mixed fields",
		input:    `{"ids": null, "authors": [], "kinds": [1]}`,
		expected: Filter{IDs: nil, Authors: []string{}, Kinds: []int{1}},
	},

	{
		name:     "zero int pointers",
		input:    `{"since": 0, "until": 0, "limit": 0}`,
		expected: Filter{Since: intPtr(0), Until: intPtr(0), Limit: intPtr(0)},
	},

	// Tags
	{
		name:     "single-letter tag",
		input:    `{"#e":["event1"]}`,
		expected: Filter{Tags: map[string][]string{"e": {"event1"}}},
	},

	{
		name:     "multi-letter tag",
		input:    `{"#emoji":["🔥","💧"]}`,
		expected: Filter{Tags: map[string][]string{"emoji": {"🔥", "💧"}}},
	},

	{
		name:     "empty tag array",
		input:    `{"#p":[]}`,
		expected: Filter{Tags: map[string][]string{"p": {}}},
	},

	{
		name:  "multiple tags",
		input: `{"#p":["pubkey1","pubkey2"],"#e":["event1","event2"]}`,
		expected: Filter{Tags: map[string][]string{
			"p": {"pubkey1", "pubkey2"},
			"e": {"event1", "event2"},
		}},
	},

	{
		name:     "null tag",
		input:    `{"#p":null}`,
		expected: Filter{Tags: map[string][]string{"p": nil}},
	},

	// Extensions
	{
		name:  "simple extension",
		input: `{"search":"query"}`,
		expected: Filter{Extensions: map[string]json.RawMessage{
			"search": json.RawMessage(`"query"`),
		},
		},
	},

	{
		name:  "extension with nested object",
		input: `{"meta":{"author":"alice","score":99}}`,
		expected: Filter{
			Extensions: map[string]json.RawMessage{
				"meta": json.RawMessage(`{"author":"alice","score":99}`),
			},
		},
	},

	{
		name:  "extension with nested array",
		input: `{"items":[1,2,3]}`,
		expected: Filter{
			Extensions: map[string]json.RawMessage{
				"items": json.RawMessage(`[1,2,3]`),
			},
		},
	},

	{
		name:  "extension with complex nested structure",
		input: `{"data":{"level1":{"level2":[{"id":1}]}}}`,
		expected: Filter{
			Extensions: map[string]json.RawMessage{
				"data": json.RawMessage(`{"level1":{"level2":[{"id":1}]}}`),
			},
		},
	},

	{
		name:  "multiple extensions",
		input: `{"search":"x","custom":true,"depth":3}`,
		expected: Filter{
			Extensions: map[string]json.RawMessage{
				"search": json.RawMessage(`"x"`),
				"custom": json.RawMessage(`true`),
				"depth":  json.RawMessage(`3`),
			},
		},
	},

	{
		name:  "extension with null value",
		input: `{"optional":null}`,
		expected: Filter{
			Extensions: map[string]json.RawMessage{
				"optional": json.RawMessage(`null`),
			},
		},
	},

	// Kitchen Sink
	{
		name:  "extension with null value",
		input: `{"ids":["x"],"since":100,"#e":["y"],"search":"z"}`,
		expected: Filter{
			IDs:   []string{"x"},
			Since: intPtr(100),
			Tags: map[string][]string{
				"e": {"y"},
			},
			Extensions: map[string]json.RawMessage{
				"search": json.RawMessage(`"z"`),
			},
		},
	},
}

var roundTripTestCases = []FilterRoundTripTestCase{
	{
		name: "fully populated filter",
		filter: Filter{
			IDs:   []string{"x"},
			Since: intPtr(100),
			Tags: map[string][]string{
				"e": {"y"},
			},
			Extensions: map[string]json.RawMessage{
				"search": json.RawMessage(`"z"`),
			},
		},
	},
}

// Tests

func TestFilterMarshalJSON(t *testing.T) {
	for _, tc := range marshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.filter.MarshalJSON()
			expectOk(t, err)

			var expectedMap, resultMap map[string]interface{}
			err = json.Unmarshal([]byte(tc.expected), &expectedMap)
			expectOk(t, err)
			err = json.Unmarshal(result, &resultMap)
			expectOk(t, err)

			if !reflect.DeepEqual(expectedMap, resultMap) {
				t.Errorf("marshal mismatch: got %s, want %s", result, tc.expected)
			}
		})
	}
}

func TestFilterUnmarshalJSON(t *testing.T) {
	for _, tc := range unmarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var result Filter
			err := result.UnmarshalJSON([]byte(tc.input))
			expectOk(t, err)

			expectEqualFilters(t, result, tc.expected)
		})
	}
}

func TestFilterRoundTrip(t *testing.T) {
	for _, tc := range roundTripTestCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBytes, err := tc.filter.MarshalJSON()
			expectOk(t, err)

			var result Filter
			err = result.UnmarshalJSON(jsonBytes)
			expectOk(t, err)

			expectEqualFilters(t, result, tc.filter)
		})
	}

}

// Helpers

func expectEqualFilters(t *testing.T, got, want Filter) {
	// Compare IDs
	if got.IDs == nil && want.IDs == nil {
		// pass
	} else if got.IDs == nil || want.IDs == nil {
		t.Errorf("mismatched ids: got %v, want %v", got.IDs, want.IDs)
	} else {
		expectEqualStringSlices(t, got.IDs, want.IDs)
	}

	// Compare Authors
	if got.Authors == nil && want.Authors == nil {
		// pass
	} else if got.Authors == nil || want.Authors == nil {
		t.Errorf("mismatched authors: got %v, want %v", got.Authors, want.Authors)
	} else {
		expectEqualStringSlices(t, got.Authors, want.Authors)
	}

	// Compare Kinds
	if got.Kinds == nil && want.Kinds == nil {
		// pass
	} else if got.Kinds == nil || want.Kinds == nil {
		t.Errorf("mismatched kinds: got %v, want %v", got.Kinds, want.Kinds)
	} else {
		expectEqualIntSlices(t, got.Kinds, want.Kinds)
	}

	// Compare Timestamps
	if got.Since == nil && want.Since == nil {
		// pass
	} else if got.Since == nil || want.Since == nil {
		t.Errorf("mismatched since pointers: got %v, want %v", got.Since, want.Since)
	} else {
		expectEqualIntPointers(t, got.Since, want.Since)
	}

	if got.Until == nil && want.Until == nil {
		// pass
	} else if got.Until == nil || want.Until == nil {
		t.Errorf("mismatched until pointers: got %v, want %v", got.Until, want.Until)
	} else {
		expectEqualIntPointers(t, got.Until, want.Until)
	}

	// Compare Limit
	if got.Limit == nil && want.Limit == nil {
		// pass
	} else if got.Limit == nil || want.Limit == nil {
		t.Errorf("mismatched limit pointers: got %v, want %v", got.Limit, want.Limit)
	} else {
		expectEqualIntPointers(t, got.Limit, want.Limit)
	}

	// Compare Tags
	if got.Tags == nil && want.Tags == nil {
		// pass
	} else if got.Tags == nil || want.Tags == nil {
		t.Errorf("mismatched tags: got %v, want %v", got.Tags, want.Tags)
	} else {
		expectEqualTags(t, got.Tags, want.Tags)
	}

	// Compare Extensions
	if got.Extensions == nil && want.Extensions == nil {
		// pass
	} else if got.Extensions == nil || want.Extensions == nil {
		t.Errorf("mismatched extensions: got %v, want %v", got.Extensions, want.Extensions)
	} else {
		expectEqualExtensions(t, got.Extensions, want.Extensions)
	}
}

func expectEqualTags(t *testing.T, got, want map[string][]string) {
	if len(got) != len(want) {
		t.Errorf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for key, wantValues := range want {
		gotValues, exists := got[key]
		if !exists {
			t.Fatalf("expected key %q to exist", key)
		}
		if len(wantValues) != len(gotValues) {
			t.Errorf(
				"key %q: length mismatch: got %d, want %d",
				key, len(gotValues), len(wantValues))
		}
		for i := range wantValues {
			if gotValues[i] != wantValues[i] {
				t.Errorf(
					"key %q: index %d: got %s, want %s",
					key, i, gotValues[i], wantValues[i])
			}
		}
	}
}

func expectEqualExtensions(t *testing.T, got, want map[string]json.RawMessage) {
	if len(got) != len(want) {
		t.Errorf("length mismatch: got %d, want %d", len(got), len(want))
	}
	for key, wantValue := range want {
		gotValue, ok := got[key]
		if !ok {
			t.Errorf("expected key %s, got nil", key)
		}
		var gotJSON, wantJSON interface{}
		if err := json.Unmarshal(wantValue, &wantJSON); err != nil {
			t.Errorf("key %q: failed to unmarshal 'want' value: %s", key, wantValue)
		}
		if err := json.Unmarshal(gotValue, &gotJSON); err != nil {
			t.Errorf("key %q: failed to unmarshal 'got' value: %s", key, wantValue)
		}
		if !reflect.DeepEqual(gotJSON, wantJSON) {
			t.Errorf("key %q: got %s, want %s", key, gotValue, wantValue)
		}
	}
}
