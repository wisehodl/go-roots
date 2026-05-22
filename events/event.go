// Roots is a purposefully minimal Nostr protocol library that provides only
// the primitives that define protocol compliance: event structure,
// serialization, cryptographic signatures, and subscription filters.
package events



// Tag represents a single tag within an event as an array of strings.
// The first element identifies the tag name, the second contains the value,
// and subsequent elements are optional.
type Tag []string

// Event represents a Nostr protocol event, with its seven required fields.
// All fields must be present for a valid event.
//easyjson:json
type Event struct {
	ID        string `json:"id"`
	PubKey    string `json:"pubkey"`
	CreatedAt int64  `json:"created_at"`
	Kind      int    `json:"kind"`
	Tags      []Tag  `json:"tags"`
	Content   string `json:"content"`
	Sig       string `json:"sig"`
}

func NewEvent(opts ...EventOption) Event {
	e := Event{}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}

type EventOption func(*Event)

func WithID(id string) EventOption {
	return func(e *Event) {
		e.ID = id
	}
}

func WithPubKey(pk string) EventOption {
	return func(e *Event) {
		e.PubKey = pk
	}
}

func WithCreatedAt(t int64) EventOption {
	return func(e *Event) {
		e.CreatedAt = t
	}
}

func WithKind(k int) EventOption {
	return func(e *Event) {
		e.Kind = k
	}
}

func WithTag(t Tag) EventOption {
	return func(e *Event) {
		e.Tags = append(e.Tags, t)
	}
}

func WithContent(c string) EventOption {
	return func(e *Event) {
		e.Content = c
	}
}

func WithSig(s string) EventOption {
	return func(e *Event) {
		e.Sig = s
	}
}
