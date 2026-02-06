package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventInterface(t *testing.T) {
	// Verify StringEvent satisfies Event interface
	var _ Event = (*StringEvent)(nil)

	concrete := &StringEvent{
		ID:        testEvent.ID,
		PubKey:    testEvent.PubKey,
		CreatedAt: testEvent.CreatedAt,
		Kind:      testEvent.Kind,
		Tags:      testEvent.Tags,
		Content:   testEvent.Content,
		Sig:       testEvent.Sig,
	}

	var event Event = concrete

	assert.Equal(t, testEvent.ID, event.GetID())
	assert.Equal(t, testEvent.PubKey, event.GetPubKey())
	assert.Equal(t, testEvent.CreatedAt, event.GetCreatedAt())
	assert.Equal(t, testEvent.Kind, event.GetKind())
	assert.Equal(t, testEvent.Tags, event.GetTags())
	assert.Equal(t, testEvent.Content, event.GetContent())
	assert.Equal(t, testEvent.Sig, event.GetSig())

	event.SetID("new-id")
	event.SetSig("new-sig")
	assert.Equal(t, "new-id", concrete.ID)
	assert.Equal(t, "new-sig", concrete.Sig)
}

func TestFunctionsAcceptInterface(t *testing.T) {
	event := &StringEvent{
		PubKey:    testEvent.PubKey,
		CreatedAt: testEvent.CreatedAt,
		Kind:      testEvent.Kind,
		Tags:      testEvent.Tags,
		Content:   testEvent.Content,
	}

	id, err := GetID(event)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	event.SetID(id)

	sig, err := SignEvent(id, testSK)
	assert.NoError(t, err)
	assert.NotEmpty(t, sig)
	event.SetSig(sig)

	err = Validate(event)
	assert.NoError(t, err)
}
