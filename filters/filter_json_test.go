package filters

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
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
		filter:   NewFilter(WithIDs(nil)),
		expected: `{}`,
	},

	{
		name:     "empty IDs",
		filter:   NewFilter(WithIDs([]string{})),
		expected: `{"ids":[]}`,
	},

	{
		name:     "populated IDs",
		filter:   NewFilter(WithIDs([]string{"abc", "123"})),
		expected: `{"ids":["abc","123"]}`,
	},

	// Author cases
	{
		name:     "nil Authors",
		filter:   NewFilter(WithAuthors(nil)),
		expected: `{}`,
	},

	{
		name:     "empty Authors",
		filter:   NewFilter(WithAuthors([]string{})),
		expected: `{"authors":[]}`,
	},

	{
		name:     "populated Authors",
		filter:   NewFilter(WithAuthors([]string{"abc", "123"})),
		expected: `{"authors":["abc","123"]}`,
	},

	// Kind cases
	{
		name:     "nil Kinds",
		filter:   NewFilter(WithKinds(nil)),
		expected: `{}`,
	},

	{
		name:     "empty Kinds",
		filter:   NewFilter(WithKinds([]int{})),
		expected: `{"kinds":[]}`,
	},

	{
		name:     "populated Kinds",
		filter:   NewFilter(WithKinds([]int{1, 20001})),
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
		filter:   NewFilter(WithSince(1000)),
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
		filter:   NewFilter(WithUntil(1000)),
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
		filter:   NewFilter(WithLimit(100)),
		expected: `{"limit":100}`,
	},

	// All standard fields
	{
		name: "all standard fields",
		filter: NewFilter(
			WithIDs([]string{"abc", "123"}),
			WithAuthors([]string{"def", "456"}),
			WithKinds([]int{1, 200, 3000}),
			WithSince(1000),
			WithUntil(2000),
			WithLimit(100),
		),
		expected: `{"ids":["abc","123"],"authors":["def","456"],"kinds":[1,200,3000],"since":1000,"until":2000,"limit":100}`,
	},

	{
		name: "mixed fields",
		filter: NewFilter(
			WithIDs(nil),
			WithAuthors([]string{}),
			WithKinds([]int{1}),
		),
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
		filter: NewFilter(
			WithTag("e", []string{"event1"}),
		),
		expected: `{"#e":["event1"]}`,
	},

	{
		name: "multi-letter tag",
		filter: NewFilter(
			WithTag("emoji", []string{"🔥", "💧"}),
		),
		expected: `{"#emoji":["🔥","💧"]}`,
	},

	{
		name: "empty tag array",
		filter: NewFilter(
			WithTag("p", []string{}),
		),
		expected: `{"#p":[]}`,
	},

	{
		name: "multiple tags",
		filter: NewFilter(
			WithTag("e", []string{"event1", "event2"}),
			WithTag("p", []string{"pubkey1", "pubkey2"}),
		),
		expected: `{"#e":["event1","event2"],"#p":["pubkey1","pubkey2"]}`,
	},

	// Extensions
	{
		name: "simple extension",
		filter: NewFilter(
			WithExtension("search", json.RawMessage(`"query"`)),
		),
		expected: `{"search":"query"}`,
	},

	{
		name: "extension with nested object",
		filter: NewFilter(
			WithExtension("meta", json.RawMessage(`{"author":"alice","score":99}`)),
		),
		expected: `{"meta":{"author":"alice","score":99}}`,
	},

	{
		name: "extension with nested array",
		filter: NewFilter(
			WithExtension("items", json.RawMessage(`[1,2,3]`)),
		),
		expected: `{"items":[1,2,3]}`,
	},

	{
		name: "extension with complex nested structure",
		filter: NewFilter(
			WithExtension("data", json.RawMessage(`{"users":[{"id":1}],"count":5}`)),
		),
		expected: `{"data":{"users":[{"id":1}],"count":5}}`,
	},

	{
		name: "multiple extensions",
		filter: NewFilter(
			WithExtension("search", json.RawMessage(`"x"`)),
			WithExtension("depth", json.RawMessage(`3`)),
		),
		expected: `{"search":"x","depth":3}`,
	},

	// Extension Collisions
	{
		name: "extension collides with standard field - IDs",
		filter: NewFilter(
			WithIDs([]string{"real"}),
			WithExtension("ids", json.RawMessage(`["fake"]`)),
		),
		expected: `{"ids":["real"]}`,
	},

	{
		name: "extension collides with standard field - Since",
		filter: NewFilter(
			WithSince(100),
			WithExtension("since", json.RawMessage(`999`)),
		),
		expected: `{"since":100}`,
	},

	{
		name: "extension collides with multiple standard fields",
		filter: NewFilter(
			WithAuthors([]string{"a"}),
			WithKinds([]int{1}),
			WithExtension("authors", json.RawMessage(`["b"]`)),
			WithExtension("kinds", json.RawMessage(`[2]`)),
		),
		expected: `{"authors":["a"],"kinds":[1]}`,
	},

	{
		name: "extension collides with tag field - #e",
		filter: NewFilter(
			WithExtension("#e", json.RawMessage(`["fakeevent"]`)),
		),
		expected: `{}`,
	},

	{
		name: "extension collides with standard and tag fields",
		filter: NewFilter(
			WithAuthors([]string{"realauthor"}),
			WithTag("e", []string{"realevent"}),
			WithExtension("authors", json.RawMessage(`["fakeauthor"]`)),
			WithExtension("#e", json.RawMessage(`["fakeevent"]`)),
		),
		expected: `{"authors":["realauthor"],"#e":["realevent"]}`,
	},

	// Kitchen Sink
	{
		name: "filter with all field types",
		filter: NewFilter(
			WithIDs([]string{"x"}),
			WithSince(100),
			WithTag("e", []string{"y"}),
			WithExtension("search", json.RawMessage(`"z"`)),
			WithExtension("ids", json.RawMessage(`["fakeid"]`)),
		),
		expected: `{"ids":["x"],"since":100,"#e":["y"],"search":"z"}`,
	},
}

var unmarshalTestCases = []FilterUnmarshalTestCase{
	{
		name:     "empty object",
		input:    `{}`,
		expected: NewFilter(),
	},

	// ID cases
	{
		name:     "null IDs",
		input:    `{"ids": null}`,
		expected: NewFilter(WithIDs(nil)),
	},

	{
		name:     "empty IDs",
		input:    `{"ids": []}`,
		expected: NewFilter(WithIDs([]string{})),
	},

	{
		name:     "populated IDs",
		input:    `{"ids": ["abc","123"]}`,
		expected: NewFilter(WithIDs([]string{"abc", "123"})),
	},

	// Author cases
	{
		name:     "null Authors",
		input:    `{"authors": null}`,
		expected: NewFilter(WithAuthors(nil)),
	},

	{
		name:     "empty Authors",
		input:    `{"authors": []}`,
		expected: NewFilter(WithAuthors([]string{})),
	},

	{
		name:     "populated Authors",
		input:    `{"authors": ["abc","123"]}`,
		expected: NewFilter(WithAuthors([]string{"abc", "123"})),
	},

	// Kind cases
	{
		name:     "null Kinds",
		input:    `{"kinds": null}`,
		expected: NewFilter(WithKinds(nil)),
	},

	{
		name:     "empty Kinds",
		input:    `{"kinds": []}`,
		expected: NewFilter(WithKinds([]int{})),
	},

	{
		name:     "populated Kinds",
		input:    `{"kinds": [1,2,3]}`,
		expected: NewFilter(WithKinds([]int{1, 2, 3})),
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
		expected: NewFilter(WithSince(1000)),
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
		expected: NewFilter(WithUntil(1000)),
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
		expected: NewFilter(WithLimit(1000)),
	},

	// All standard fields
	{
		name:  "all standard fields",
		input: `{"ids":["abc","123"],"authors":["def","456"],"kinds":[1,200,3000],"since":1000,"until":2000,"limit":100}`,
		expected: NewFilter(
			WithIDs([]string{"abc", "123"}),
			WithAuthors([]string{"def", "456"}),
			WithKinds([]int{1, 200, 3000}),
			WithSince(1000),
			WithUntil(2000),
			WithLimit(100),
		),
	},

	{
		name:  "mixed fields",
		input: `{"ids": null, "authors": [], "kinds": [1]}`,
		expected: NewFilter(
			WithIDs(nil),
			WithAuthors([]string{}),
			WithKinds([]int{1}),
		),
	},

	{
		name:     "zero int pointers",
		input:    `{"since": 0, "until": 0, "limit": 0}`,
		expected: NewFilter(WithSince(0), WithUntil(0), WithLimit(0)),
	},

	// Tags
	{
		name:     "single-letter tag",
		input:    `{"#e":["event1"]}`,
		expected: NewFilter(WithTag("e", []string{"event1"})),
	},

	{
		name:     "multi-letter tag",
		input:    `{"#emoji":["🔥","💧"]}`,
		expected: NewFilter(WithTag("emoji", []string{"🔥", "💧"})),
	},

	{
		name:     "empty tag array",
		input:    `{"#p":[]}`,
		expected: NewFilter(WithTag("p", []string{})),
	},

	{
		name:  "multiple tags",
		input: `{"#p":["pubkey1","pubkey2"],"#e":["event1","event2"]}`,
		expected: NewFilter(
			WithTag("p", []string{"pubkey1", "pubkey2"}),
			WithTag("e", []string{"event1", "event2"}),
		),
	},

	{
		name:     "null tag",
		input:    `{"#p":null}`,
		expected: NewFilter(WithTag("p", nil)),
	},

	// Extensions
	{
		name:  "simple extension",
		input: `{"search":"query"}`,
		expected: NewFilter(
			WithExtension("search", json.RawMessage(`"query"`)),
		),
	},

	{
		name:  "extension with nested object",
		input: `{"meta":{"author":"alice","score":99}}`,
		expected: NewFilter(
			WithExtension("meta", json.RawMessage(`{"author":"alice","score":99}`)),
		),
	},

	{
		name:  "extension with nested array",
		input: `{"items":[1,2,3]}`,
		expected: NewFilter(
			WithExtension("items", json.RawMessage(`[1,2,3]`)),
		),
	},

	{
		name:  "extension with complex nested structure",
		input: `{"data":{"level1":{"level2":[{"id":1}]}}}`,
		expected: NewFilter(
			WithExtension("data", json.RawMessage(`{"level1":{"level2":[{"id":1}]}}`)),
		),
	},

	{
		name:  "multiple extensions",
		input: `{"search":"x","custom":true,"depth":3}`,
		expected: NewFilter(
			WithExtension("search", json.RawMessage(`"x"`)),
			WithExtension("custom", json.RawMessage(`true`)),
			WithExtension("depth", json.RawMessage(`3`)),
		),
	},

	{
		name:  "extension with null value",
		input: `{"optional":null}`,
		expected: NewFilter(
			WithExtension("optional", json.RawMessage(`null`)),
		),
	},

	// Kitchen Sink
	{
		name:  "extension with null value",
		input: `{"ids":["x"],"since":100,"#e":["y"],"search":"z"}`,
		expected: NewFilter(
			WithIDs([]string{"x"}),
			WithSince(100),
			WithTag("e", []string{"y"}),
			WithExtension("search", json.RawMessage(`"z"`)),
		),
	},
}

var roundTripTestCases = []FilterRoundTripTestCase{
	{
		name: "fully populated filter",
		filter: NewFilter(
			WithIDs([]string{"x"}),
			WithSince(100),
			WithTag("e", []string{"y"}),
			WithExtension("search", json.RawMessage(`"z"`)),
		),
	},
}

// Tests

func TestFilterMarshalJSON(t *testing.T) {
	for _, tc := range marshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MarshalJSON(tc.filter)
			assert.NoError(t, err)

			var expectedMap, actualMap map[string]interface{}
			err = json.Unmarshal([]byte(tc.expected), &expectedMap)
			assert.NoError(t, err)
			err = json.Unmarshal(result, &actualMap)
			assert.NoError(t, err)

			assert.Equal(t, expectedMap, actualMap)
		})
	}
}

func TestFilterUnmarshalJSON(t *testing.T) {
	for _, tc := range unmarshalTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var result Filter
			err := UnmarshalJSON([]byte(tc.input), &result)
			assert.NoError(t, err)

			expectEqualFilters(t, result, tc.expected)
		})
	}
}

func TestFilterRoundTrip(t *testing.T) {
	for _, tc := range roundTripTestCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonBytes, err := MarshalJSON(tc.filter)
			assert.NoError(t, err)

			var result Filter
			err = UnmarshalJSON(jsonBytes, &result)
			assert.NoError(t, err)

			expectEqualFilters(t, result, tc.filter)
		})
	}

}

// Helpers

func expectEqualFilters(t *testing.T, got, want Filter) {
	assert.Equal(t, want.IDs, got.IDs)
	assert.Equal(t, want.Authors, got.Authors)
	assert.Equal(t, want.Kinds, got.Kinds)
	assert.Equal(t, want.Since, got.Since)
	assert.Equal(t, want.Until, got.Until)
	assert.Equal(t, want.Limit, got.Limit)
	assert.Equal(t, want.Tags, got.Tags)

	if want.Extensions == nil && got.Extensions == nil {
		return
	}
	assert.NotNil(t, got.Extensions)
	assert.NotNil(t, want.Extensions)

	assert.Equal(t, len(want.Extensions), len(got.Extensions))
	for key, wantValue := range want.Extensions {
		gotValue, ok := got.Extensions[key]
		assert.True(t, ok, "expected key %s", key)

		var gotJSON, wantJSON interface{}
		assert.NoError(t, json.Unmarshal(wantValue, &wantJSON))
		assert.NoError(t, json.Unmarshal(gotValue, &gotJSON))
		assert.Equal(t, wantJSON, gotJSON)
	}

}
