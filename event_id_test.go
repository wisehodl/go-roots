package roots

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type IDTestCase struct {
	name     string
	event    Event
	expected string
}

var idTestCases = []IDTestCase{
	{
		name: "minimal event",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags:      [][]string{},
			Content:   "",
		},
		expected: "13a55672a600398894592f4cb338652d4936caffe5d3718d11597582bb030c39",
	},

	{
		name: "alphanumeric content",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags:      [][]string{},
			Content:   "hello world",
		},
		expected: "c7a702e6158744ca03508bbb4c90f9dbb0d6e88fefbfaa511d5ab24b4e3c48ad",
	},

	{
		name: "unicode content",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags:      [][]string{},
			Content:   "hello world 😀",
		},
		expected: "e42083fafbf9a39f97914fd9a27cedb38c429ac3ca8814288414eaad1f472fe8",
	},

	{
		name: "escaped content",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags:      [][]string{},
			Content:   "\"You say yes.\"\\n\\t\"I say no.\"",
		},
		expected: "343de133996a766bf00561945b6f2b2717d4905275976ca75c1d7096b7d1900c",
	},

	{
		name: "json content",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags:      [][]string{},
			Content:   "{\"field\": [\"value\",\"value\"],\"numeral\": 123}",
		},
		expected: "c6140190453ee947efb790e70541a9d37c41604d1f29e4185da4325621ed5270",
	},

	{
		name: "empty tag",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags: [][]string{
				{"a", ""},
			},
			Content: "",
		},
		expected: "7d3e394c75916362436f11c603b1a89b40b50817550cfe522a90d769655007a4",
	},

	{
		name: "single tag",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags: [][]string{
				{"a", "value"},
			},
			Content: "",
		},
		expected: "7db394e274fb893edbd9f4aa9ff189d4f3264bf1a29cef8f614e83ebf6fa19fe",
	},

	{
		name: "optional tag values",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags: [][]string{
				{"a", "value", "optional"},
			},
			Content: "",
		},
		expected: "656b47884200959e0c03054292c453cfc4beea00b592d92c0f557bff765e9d34",
	},

	{
		name: "multiple tags",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags: [][]string{
				{"a", "value", "optional"},
				{"b", "another"},
				{"c", "data"},
			},
			Content: "",
		},
		expected: "f7c27f2eacda7ece5123a4f82db56145ba59f7c9e6c5eeb88552763664506b06",
	},

	{
		name: "unicode tag",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      1,
			Tags: [][]string{
				{"a", "😀"},
			},
			Content: "",
		},
		expected: "fd2798d165d9bf46acbe817735dc8cedacd4c42dfd9380792487d4902539e986",
	},

	{
		name: "zero timestamp",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: 0,
			Kind:      1,
			Tags:      [][]string{},
			Content:   "",
		},
		expected: "9ca742f2e2eea72ad6e0277a6287e2bb16a3e47d64b8468bc98474e266cf0ec2",
	},

	{
		name: "negative timestamp",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: -1760740551,
			Kind:      1,
			Tags:      [][]string{},
			Content:   "",
		},
		expected: "4740b027040bb4d0ee8e885f567a80277097da70cddd143d8a6dadf97f6faaa3",
	},

	{
		name: "max int64 timestamp",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: 9223372036854775807,
			Kind:      1,
			Tags:      [][]string{},
			Content:   "",
		},
		expected: "b28cdd44496acb49e36c25859f0f819122829a12dc57c07612d5f44cb121d2a7",
	},

	{
		name: "different kind",
		event: Event{
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      20021,
			Tags:      [][]string{},
			Content:   "",
		},
		expected: "995c4894c264e6b9558cb94b7b34008768d53801b99960b47298d4e3e23fadd3",
	},
}

func TestEventGetId(t *testing.T) {
	for _, tc := range idTestCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := tc.event.GetID()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
