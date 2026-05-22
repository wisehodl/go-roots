package events

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testEventWithTags = NewEvent(
	WithID("c9a0f84fcaa889654da8992105eb122eb210c8cbd58210609a5ef7e170b51400"),
	WithPubKey(testPK),
	WithCreatedAt(1760740551),
	WithKind(1),
	WithTag(Tag{"a", "value"}),
	WithTag(Tag{"b", "value", "optional"}),
	WithContent("valid event"),
	WithSig("668a715f1eb983172acf230d17bd283daedb2598adf8de4290bcc7eb0b802fdb60669d1e7d1104ac70393f4dbccd07e8abf897152af6ce6c0a75499874e27f14"),
)

func TestValidatedEventConstruction(t *testing.T) {
	t.Run("accepts valid event", func(t *testing.T) {
		_, err := NewValidatedEvent(testEvent)
		assert.NoError(t, err)
	})

	t.Run("rejects invalid structure", func(t *testing.T) {
		bad := NewEvent(
			WithID(testEvent.ID),
			WithPubKey("notavalidpubkey"),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		)
		_, err := NewValidatedEvent(bad)
		assert.ErrorContains(t, err, "public key must be 64 lowercase hex characters")
	})

	t.Run("rejects id mismatch", func(t *testing.T) {
		bad := NewEvent(
			WithID("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			WithPubKey(testEvent.PubKey),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		)
		_, err := NewValidatedEvent(bad)
		assert.ErrorContains(t, err, "does not match computed id")
	})

	t.Run("rejects invalid signature", func(t *testing.T) {
		bad := NewEvent(
			WithID(testEvent.ID),
			WithPubKey(testEvent.PubKey),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		)
		_, err := NewValidatedEvent(bad)
		assert.ErrorContains(t, err, "event signature is invalid")
	})
}

func TestValidatedEventAccessors(t *testing.T) {
	ve, err := NewValidatedEvent(testEventWithTags)
	assert.NoError(t, err)

	t.Run("ID", func(t *testing.T) {
		assert.Equal(t, testEventWithTags.ID, ve.ID())
	})

	t.Run("PubKey", func(t *testing.T) {
		assert.Equal(t, testEventWithTags.PubKey, ve.PubKey())
	})

	t.Run("CreatedAt", func(t *testing.T) {
		assert.Equal(t, testEventWithTags.CreatedAt, ve.CreatedAt())
	})

	t.Run("Kind", func(t *testing.T) {
		assert.Equal(t, testEventWithTags.Kind, ve.Kind())
	})

	t.Run("Content", func(t *testing.T) {
		assert.Equal(t, testEventWithTags.Content, ve.Content())
	})

	t.Run("Sig", func(t *testing.T) {
		assert.Equal(t, testEventWithTags.Sig, ve.Sig())
	})

	t.Run("Tags", func(t *testing.T) {
		assert.Equal(t, testEventWithTags.Tags, ve.Tags())
	})
}

func TestValidatedEventTagsImmutability(t *testing.T) {
	t.Run("outer slice mutation is isolated", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEventWithTags)
		tags := ve.Tags()
		tags[0] = Tag{"z", "injected"}
		assert.Equal(t, Tag{"a", "value"}, ve.Tags()[0])
	})

	t.Run("inner slice mutation is isolated", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEventWithTags)
		tags := ve.Tags()
		tags[0][1] = "mutated"
		assert.Equal(t, "value", ve.Tags()[0][1])
	})

	t.Run("append to outer slice is isolated", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEventWithTags)
		originalLen := len(ve.Tags())
		tags := ve.Tags()
		tags = append(tags, Tag{"new", "tag"})
		_ = tags
		assert.Equal(t, originalLen, len(ve.Tags()))
	})

	t.Run("append to inner slice is isolated", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEventWithTags)
		originalLen := len(ve.Tags()[1])
		tags := ve.Tags()
		tags[1] = append(tags[1], "extra")
		assert.Equal(t, originalLen, len(ve.Tags()[1]))
	})

	t.Run("successive calls return independent copies", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEventWithTags)
		first := ve.Tags()
		second := ve.Tags()
		first[0][0] = "z"
		assert.Equal(t, "a", second[0][0])
	})

	t.Run("nil tags remain nil", func(t *testing.T) {
		nil_tags_event := Event{
			ID:        testEvent.ID,
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      nil,
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		}
		ve, err := NewValidatedEvent(nil_tags_event)
		assert.NoError(t, err)
		assert.Nil(t, ve.Tags())
		assert.Nil(t, ve.Event().Tags)
	})
}

func TestValidatedEventEventImmutability(t *testing.T) {
	t.Run("outer tags mutation is isolated", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEventWithTags)
		e := ve.Event()
		e.Tags[0] = Tag{"z", "injected"}
		assert.Equal(t, Tag{"a", "value"}, ve.Tags()[0])
	})

	t.Run("inner tags mutation is isolated", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEventWithTags)
		e := ve.Event()
		e.Tags[0][1] = "mutated"
		assert.Equal(t, "value", ve.Tags()[0][1])
	})

	t.Run("scalar field mutation is isolated", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEvent)
		e := ve.Event()
		e.Content = "tampered"
		assert.Equal(t, testEvent.Content, ve.Content())
	})

	t.Run("successive calls return independent copies", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEventWithTags)
		first := ve.Event()
		second := ve.Event()
		first.Tags[0][0] = "z"
		assert.Equal(t, "a", second.Tags[0][0])
	})

	t.Run("returned event passes validation", func(t *testing.T) {
		ve, _ := NewValidatedEvent(testEventWithTags)
		assert.NoError(t, Validate(ve.Event()))
	})

	t.Run("source event mutation after construction is isolated", func(t *testing.T) {
		source := NewEvent(
			WithID(testEventWithTags.ID),
			WithPubKey(testEventWithTags.PubKey),
			WithCreatedAt(testEventWithTags.CreatedAt),
			WithKind(testEventWithTags.Kind),
			WithTag(Tag{"a", "value"}),
			WithTag(Tag{"b", "value", "optional"}),
			WithContent(testEventWithTags.Content),
			WithSig(testEventWithTags.Sig),
		)
		ve, err := NewValidatedEvent(source)
		assert.NoError(t, err)
		source.Tags[0][1] = "mutated"
		source.Content = "tampered"
		assert.Equal(t, "value", ve.Tags()[0][1])
		assert.Equal(t, testEventWithTags.Content, ve.Content())
	})
}

func TestValidatedEventMarshalJSON(t *testing.T) {
	ve, _ := NewValidatedEvent(testEvent)
	jsonBytes, err := json.Marshal(ve)
	assert.NoError(t, err)
	assert.Equal(t, testEventJSON, string(jsonBytes))
}
