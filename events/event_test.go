package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testSK = "f43a0435f69529f310bbd1d6263d2fbf0977f54bfe2310cc37ae5904b83bb167"
const testPK = "cfa87f35acbde29ba1ab3ee42de527b2cad33ac487e80cf2d6405ea0042c8fef"

var testEvent = Event{
	ID:        "c7a702e6158744ca03508bbb4c90f9dbb0d6e88fefbfaa511d5ab24b4e3c48ad",
	PubKey:    testPK,
	CreatedAt: 1760740551,
	Kind:      1,
	Tags:      []Tag{},
	Content:   "hello world",
	Sig:       "83b71e15649c9e9da362c175f988c36404cabf357a976d869102a74451cfb8af486f6088b5631033b4927bd46cad7a0d90d7f624aefc0ac260364aa65c36071a",
}

var testEventJSON = `{"id":"c7a702e6158744ca03508bbb4c90f9dbb0d6e88fefbfaa511d5ab24b4e3c48ad","pubkey":"cfa87f35acbde29ba1ab3ee42de527b2cad33ac487e80cf2d6405ea0042c8fef","created_at":1760740551,"kind":1,"tags":[],"content":"hello world","sig":"83b71e15649c9e9da362c175f988c36404cabf357a976d869102a74451cfb8af486f6088b5631033b4927bd46cad7a0d90d7f624aefc0ac260364aa65c36071a"}`
var testEventJSONBytes = []byte(testEventJSON)

func expectEqualEvents(t *testing.T, got, want Event) {
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.PubKey, got.PubKey)
	assert.Equal(t, want.CreatedAt, got.CreatedAt)
	assert.Equal(t, want.Kind, got.Kind)
	assert.Equal(t, want.Content, got.Content)
	assert.Equal(t, want.Sig, got.Sig)
	assert.Equal(t, want.Tags, got.Tags)
}
