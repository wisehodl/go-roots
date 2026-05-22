package filters

import (
	"encoding/json"
	"git.wisehodl.dev/jay/go-roots/events"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testEvents []events.ValidatedEvent

func init() {
	data, err := os.ReadFile("testdata/test_events.json")
	if err != nil {
		panic(err)
	}
	var raw []events.Event
	if err := json.Unmarshal(data, &raw); err != nil {
		panic(err)
	}
	for _, e := range raw {
		ve, err := events.NewValidatedEvent(e)
		if err != nil {
			panic(err)
		}
		testEvents = append(testEvents, ve)
	}
}

// Test keypairs corresponding to test events, for reference.
const (
	nayru_sk  = "1784be782585dfa97712afe12585d13ee608b624cf564116fa143c31a124d31e"
	nayru_pk  = "d877e187934bd942a71221b50ff2b426bd0777991b41b6c749119805dc40bcbe"
	farore_sk = "03d0611c41048a9108a75bf5d023180b5cf2d2d24e2e6b83def29de977315bb3"
	farore_pk = "9e4b726ab0f25af580bdd2fd504fb245cf604f1fbc2482b89cf74beb4fb3aca9"
	din_sk    = "7547dd630c04fde72bff3b99c481c683479966cb758f0b367b08fc971ead18f0"
	din_pk    = "e719e8f83b77a9efacb29fd19118b030cbf7cfbca1f8d3694235707ee213abc7"
)

type FilterTestCase struct {
	name        string
	filter      Filter
	expectedIDs []string
}

var filterTestCases = []FilterTestCase{
	{
		name:   "empty filter",
		filter: NewFilter(),
		expectedIDs: []string{
			"e751d41f",
			"2e06c187",
			"e67fa7b8",
			"5e4c64f1",
			"7a5d83d4",
			"3a122100",
			"4a15d963",
			"4b03b69a",
			"d39e6f3f",
		},
	},

	{
		name:   "empty id",
		filter: NewFilter(WithIDs([]string{})),
		expectedIDs: []string{
			"e751d41f",
			"2e06c187",
			"e67fa7b8",
			"5e4c64f1",
			"7a5d83d4",
			"3a122100",
			"4a15d963",
			"4b03b69a",
			"d39e6f3f",
		},
	},

	{
		name:        "single id prefix",
		filter:      NewFilter(WithIDs([]string{"e751d41f"})),
		expectedIDs: []string{"e751d41f"},
	},

	{
		name: "single full id",
		filter: NewFilter(
			WithIDs([]string{
				"e67fa7b84df6b0bb4c57f8719149de77f58955d7849da1be10b2267c72daad8b"}),
		),
		expectedIDs: []string{"e67fa7b8"},
	},

	{
		name:        "multiple id prefixes",
		filter:      NewFilter(WithIDs([]string{"2e06c187", "5e4c64f1"})),
		expectedIDs: []string{"2e06c187", "5e4c64f1"},
	},

	{
		name:        "no id match",
		filter:      NewFilter(WithIDs([]string{"ffff"})),
		expectedIDs: []string{},
	},

	{
		name:   "empty author",
		filter: NewFilter(WithAuthors([]string{})),
		expectedIDs: []string{
			"e751d41f",
			"2e06c187",
			"e67fa7b8",
			"5e4c64f1",
			"7a5d83d4",
			"3a122100",
			"4a15d963",
			"4b03b69a",
			"d39e6f3f",
		},
	},

	{
		name:        "single author prefix",
		filter:      NewFilter(WithAuthors([]string{"d877e187"})),
		expectedIDs: []string{"e751d41f", "2e06c187", "e67fa7b8"},
	},

	{
		name:   "multiple author prefixex",
		filter: NewFilter(WithAuthors([]string{"d877e187", "9e4b726a"})),
		expectedIDs: []string{
			"e751d41f",
			"2e06c187",
			"e67fa7b8",
			"5e4c64f1",
			"7a5d83d4",
			"3a122100",
		},
	},

	{
		name: "single author full",
		filter: NewFilter(
			WithAuthors([]string{
				"d877e187934bd942a71221b50ff2b426bd0777991b41b6c749119805dc40bcbe"}),
		),
		expectedIDs: []string{"e751d41f", "2e06c187", "e67fa7b8"},
	},

	{
		name:        "no author match",
		filter:      NewFilter(WithAuthors([]string{"ffff"})),
		expectedIDs: []string{},
	},

	{
		name:   "empty kind",
		filter: NewFilter(WithKinds([]int{})),
		expectedIDs: []string{
			"e751d41f",
			"2e06c187",
			"e67fa7b8",
			"5e4c64f1",
			"7a5d83d4",
			"3a122100",
			"4a15d963",
			"4b03b69a",
			"d39e6f3f",
		},
	},

	{
		name:        "single kind",
		filter:      NewFilter(WithKinds([]int{1})),
		expectedIDs: []string{"2e06c187", "7a5d83d4", "4b03b69a"},
	},

	{
		name:   "multiple kinds",
		filter: NewFilter(WithKinds([]int{0, 2})),
		expectedIDs: []string{
			"e751d41f",
			"e67fa7b8",
			"5e4c64f1",
			"3a122100",
			"4a15d963",
			"d39e6f3f",
		},
	},

	{
		name:        "no kind match",
		filter:      NewFilter(WithKinds([]int{99})),
		expectedIDs: []string{},
	},

	{
		name:   "since only",
		filter: NewFilter(WithSince(5000)),
		expectedIDs: []string{
			"7a5d83d4",
			"3a122100",
			"4a15d963",
			"4b03b69a",
			"d39e6f3f",
		},
	},

	{
		name:   "until only",
		filter: NewFilter(WithUntil(3000)),
		expectedIDs: []string{
			"e751d41f",
			"2e06c187",
			"e67fa7b8",
		},
	},

	{
		name:   "time range",
		filter: NewFilter(WithSince(4000), WithUntil(6000)),
		expectedIDs: []string{
			"5e4c64f1",
			"7a5d83d4",
			"3a122100",
		},
	},

	{
		name:        "outside time range",
		filter:      NewFilter(WithSince(10000)),
		expectedIDs: []string{},
	},

	{
		name:   "empty tag filter",
		filter: NewFilter(WithTag("e", []string{})),
		expectedIDs: []string{
			"e751d41f",
			"2e06c187",
			"e67fa7b8",
			"5e4c64f1",
			"7a5d83d4",
			"3a122100",
			"4a15d963",
			"4b03b69a",
			"d39e6f3f",
		},
	},

	{
		name: "single letter tag filter: e",
		filter: NewFilter(
			WithTag("e", []string{
				"5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36"}),
		),
		expectedIDs: []string{"2e06c187"},
	},

	{
		name: "multiple tag matches",
		filter: NewFilter(
			WithTag("e", []string{
				"5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36",
				"ae3f2a91b6c3d8f7e9a1c5b4d8f2e7a9b6c3d8f7e9a1c5b4d8f2e7a9b6c3d8f7",
			}),
		),
		expectedIDs: []string{"2e06c187", "3a122100"},
	},

	{
		name: "multiple tag matches - single event match",
		filter: NewFilter(
			WithTag("e", []string{
				"5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36",
				"cb7787c460a79187d6a13e75a0f19240e05fafca8ea42288f5765773ea69cf2f",
			}),
		),
		expectedIDs: []string{"2e06c187"},
	},

	{
		name: "single letter tag filter: p",
		filter: NewFilter(
			WithTag("p", []string{
				"91cf9b32f3735070f46c0a86a820a47efa08a5be6c9f4f8cf68e5b5b75c92d60"}),
		),
		expectedIDs: []string{"e67fa7b8"},
	},

	{
		name: "multi letter tag filter",
		filter: NewFilter(
			WithTag("emoji", []string{"🌊"}),
		),
		expectedIDs: []string{"e67fa7b8"},
	},

	{
		name: "multiple tag filters",
		filter: NewFilter(
			WithTag("e", []string{
				"ae3f2a91b6c3d8f7e9a1c5b4d8f2e7a9b6c3d8f7e9a1c5b4d8f2e7a9b6c3d8f7"}),
			WithTag("p", []string{
				"3bf0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d"}),
		),
		expectedIDs: []string{"3a122100"},
	},

	{
		name: "prefix tag filter",
		filter: NewFilter(
			WithTag("p", []string{"ae3f2a91"}),
		),
		expectedIDs: []string{},
	},

	{
		name: "unknown tag filter",
		filter: NewFilter(
			WithTag("z", []string{"anything"}),
		),
		expectedIDs: []string{},
	},

	{
		name: "combined author+kind tag filter",
		filter: NewFilter(
			WithAuthors([]string{"d877e187"}),
			WithKinds([]int{1, 2}),
		),
		expectedIDs: []string{
			"2e06c187",
			"e67fa7b8",
		},
	},

	{
		name: "combined kind+time range tag filter",
		filter: NewFilter(
			WithKinds([]int{0}),
			WithSince(2000),
			WithUntil(7000),
		),
		expectedIDs: []string{
			"5e4c64f1",
			"4a15d963",
		},
	},

	{
		name: "combined author+tag tag filter",
		filter: NewFilter(
			WithAuthors([]string{"e719e8f8"}),
			WithTag("power", []string{"fire"}),
		),
		expectedIDs: []string{
			"4a15d963",
		},
	},

	{
		name: "combined tag filter",
		filter: NewFilter(
			WithAuthors([]string{"e719e8f8"}),
			WithKinds([]int{0}),
			WithSince(5000),
			WithUntil(10000),
			WithTag("power", []string{"fire"}),
		),
		expectedIDs: []string{
			"4a15d963",
		},
	},
}

func TestEventFilterMatching(t *testing.T) {
	for _, tc := range filterTestCases {
		t.Run(tc.name, func(t *testing.T) {
			actualIDs := []string{}
			for _, event := range testEvents {
				if Matches(tc.filter, event) {
					actualIDs = append(actualIDs, event.ID()[:8])
				}
			}

			assert.Equal(t, tc.expectedIDs, actualIDs)
		})
	}
}
