// Roots is a purposefully minimal Nostr protocol library that provides only
// the primitives that define protocol compliance: event structure,
// serialization, cryptographic signatures, and subscription filters.
package events

import (
	"regexp"
)

// Tag represents a single tag within an event as an array of strings.
// The first element identifies the tag name, the second contains the value,
// and subsequent elements are optional.
type Tag []string

// Event is the stable contract for Nostr event operations.
// Implementations may optimize internal representation while
// satisfying this interface.
type Event interface {
	GetID() string
	GetPubKey() string
	GetCreatedAt() int
	GetKind() int
	GetTags() []Tag
	GetContent() string
	GetSig() string

	// SetID sets the event ID after computation
	SetID(id string)

	// SetSig sets the signature after signing
	SetSig(sig string)
}

var (
	// Hex64Pattern matches 64-character, lowercase, hexadecimal strings.
	// Used for validating event IDs and cryptographic keys.
	Hex64Pattern = regexp.MustCompile("^[a-f0-9]{64}$")

	// Hex128Pattern matches 128-character, lowercase, hexadecimal strings.
	// Used for validating signatures.
	Hex128Pattern = regexp.MustCompile("^[a-f0-9]{128}$")
)
