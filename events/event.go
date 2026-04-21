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
type Event struct {
	ID        string `json:"id"`
	PubKey    string `json:"pubkey"`
	CreatedAt int    `json:"created_at"`
	Kind      int    `json:"kind"`
	Tags      []Tag  `json:"tags"`
	Content   string `json:"content"`
	Sig       string `json:"sig"`
}
