package roots

import (
	"encoding/json"
	"os"
	"testing"
)

var testEvents []Event

func init() {
	data, err := os.ReadFile("testdata/test_events.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &testEvents); err != nil {
		panic(err)
	}
}

var (
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
	matchingIDs []string
}

var filterTestCases = []FilterTestCase{
	{
		name:   "empty filter",
		filter: Filter{},
		matchingIDs: []string{
			"e751d41f",
			"562bc378",
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
		filter: Filter{IDs: []string{}},
		matchingIDs: []string{
			"e751d41f",
			"562bc378",
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
		filter:      Filter{IDs: []string{"e751d41f"}},
		matchingIDs: []string{"e751d41f"},
	},

	{
		name:        "single full id",
		filter:      Filter{IDs: []string{"e67fa7b84df6b0bb4c57f8719149de77f58955d7849da1be10b2267c72daad8b"}},
		matchingIDs: []string{"e67fa7b8"},
	},

	{
		name:        "multiple id prefixes",
		filter:      Filter{IDs: []string{"562bc378", "5e4c64f1"}},
		matchingIDs: []string{"562bc378", "5e4c64f1"},
	},

	{
		name:        "no id match",
		filter:      Filter{IDs: []string{"ffff"}},
		matchingIDs: []string{},
	},

	{
		name:   "empty author",
		filter: Filter{Authors: []string{}},
		matchingIDs: []string{
			"e751d41f",
			"562bc378",
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
		filter:      Filter{Authors: []string{"d877e187"}},
		matchingIDs: []string{"e751d41f", "562bc378", "e67fa7b8"},
	},

	{
		name:   "multiple author prefixex",
		filter: Filter{Authors: []string{"d877e187", "9e4b726a"}},
		matchingIDs: []string{
			"e751d41f",
			"562bc378",
			"e67fa7b8",
			"5e4c64f1",
			"7a5d83d4",
			"3a122100",
		},
	},

	{
		name:        "single author full",
		filter:      Filter{Authors: []string{"d877e187934bd942a71221b50ff2b426bd0777991b41b6c749119805dc40bcbe"}},
		matchingIDs: []string{"e751d41f", "562bc378", "e67fa7b8"},
	},

	{
		name:        "no author match",
		filter:      Filter{Authors: []string{"ffff"}},
		matchingIDs: []string{},
	},

	{
		name:   "empty kind",
		filter: Filter{Kinds: []int{}},
		matchingIDs: []string{
			"e751d41f",
			"562bc378",
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
		filter:      Filter{Kinds: []int{1}},
		matchingIDs: []string{"562bc378", "7a5d83d4", "4b03b69a"},
	},

	{
		name:   "multiple kinds",
		filter: Filter{Kinds: []int{0, 2}},
		matchingIDs: []string{
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
		filter:      Filter{Kinds: []int{99}},
		matchingIDs: []string{},
	},

	{
		name:   "since only",
		filter: Filter{Since: intPtr(5000)},
		matchingIDs: []string{
			"7a5d83d4",
			"3a122100",
			"4a15d963",
			"4b03b69a",
			"d39e6f3f",
		},
	},

	{
		name:   "until only",
		filter: Filter{Until: intPtr(3000)},
		matchingIDs: []string{
			"e751d41f",
			"562bc378",
			"e67fa7b8",
		},
	},

	{
		name: "time range",
		filter: Filter{
			Since: intPtr(4000),
			Until: intPtr(6000),
		},
		matchingIDs: []string{
			"5e4c64f1",
			"7a5d83d4",
			"3a122100",
		},
	},

	{
		name: "outside time range",
		filter: Filter{
			Since: intPtr(10000),
		},
		matchingIDs: []string{},
	},

	{
		name: "empty tag filter",
		filter: Filter{
			Tags: map[string][]string{
				"e": {},
			},
		},
		matchingIDs: []string{
			"e751d41f",
			"562bc378",
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
		filter: Filter{
			Tags: map[string][]string{
				"e": {"5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36"},
			},
		},
		matchingIDs: []string{"562bc378"},
	},

	{
		name: "multiple tag matches",
		filter: Filter{
			Tags: map[string][]string{
				"e": {
					"5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36",
					"ae3f2a91b6c3d8f7e9a1c5b4d8f2e7a9b6c3d8f7e9a1c5b4d8f2e7a9b6c3d8f7",
				},
			},
		},
		matchingIDs: []string{"562bc378", "3a122100"},
	},

	{
		name: "multiple tag matches - single event match",
		filter: Filter{
			Tags: map[string][]string{
				"e": {
					"5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36",
					"cb7787c460a79187d6a13e75a0f19240e05fafca8ea42288f5765773ea69cf2f",
				},
			},
		},
		matchingIDs: []string{"562bc378"},
	},

	{
		name: "single letter tag filter: p",
		filter: Filter{
			Tags: map[string][]string{
				"p": {"91cf9b32f3735070f46c0a86a820a47efa08a5be6c9f4f8cf68e5b5b75c92d60"},
			},
		},
		matchingIDs: []string{"e67fa7b8"},
	},

	{
		name: "multi letter tag filter",
		filter: Filter{
			Tags: map[string][]string{
				"emoji": {"🌊"},
			},
		},
		matchingIDs: []string{"e67fa7b8"},
	},

	{
		name: "multiple tag filters",
		filter: Filter{
			Tags: map[string][]string{
				"e": {"ae3f2a91b6c3d8f7e9a1c5b4d8f2e7a9b6c3d8f7e9a1c5b4d8f2e7a9b6c3d8f7"},
				"p": {"3bf0c63fcb93463407af97a5e5ee64fa883d107ef9e558472c4eb9aaaefa459d"},
			},
		},
		matchingIDs: []string{"3a122100"},
	},

	{
		name: "prefix tag filter",
		filter: Filter{
			Tags: map[string][]string{
				"p": {"ae3f2a91"},
			},
		},
		matchingIDs: []string{},
	},

	{
		name: "unknown tag filter",
		filter: Filter{
			Tags: map[string][]string{
				"z": {"anything"},
			},
		},
		matchingIDs: []string{},
	},

	{
		name: "combined author+kind tag filter",
		filter: Filter{
			Authors: []string{"d877e187"},
			Kinds:   []int{1, 2},
		},
		matchingIDs: []string{
			"562bc378",
			"e67fa7b8",
		},
	},

	{
		name: "combined kind+time range tag filter",
		filter: Filter{
			Kinds: []int{0},
			Since: intPtr(2000),
			Until: intPtr(7000),
		},
		matchingIDs: []string{
			"5e4c64f1",
			"4a15d963",
		},
	},

	{
		name: "combined author+tag tag filter",
		filter: Filter{
			Authors: []string{"e719e8f8"},
			Tags: map[string][]string{
				"power": {"fire"},
			},
		},
		matchingIDs: []string{
			"4a15d963",
		},
	},

	{
		name: "combined tag filter",
		filter: Filter{
			Authors: []string{"e719e8f8"},
			Kinds:   []int{0},
			Since:   intPtr(5000),
			Until:   intPtr(10000),
			Tags: map[string][]string{
				"power": {"fire"},
			},
		},
		matchingIDs: []string{
			"4a15d963",
		},
	},
}

func TestEventFilterMatching(t *testing.T) {
	for _, tc := range filterTestCases {
		t.Run(tc.name, func(t *testing.T) {
			matchedIDs := []string{}
			for _, event := range testEvents {
				if tc.filter.Matches(event) {
					matchedIDs = append(matchedIDs, event.ID[:8])
				}
			}

			expectEqualStringSlices(t, matchedIDs, tc.matchingIDs)
		})
	}
}
