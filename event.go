// Roots is a purposefully minimal Nostr protocol library that provides only
// the primitives that define protocol compliance: event structure,
// serialization, cryptographic signatures, and subscription filters.
package roots

import (
	"errors"
	"regexp"
)

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

var (
	// Hex64Pattern matches 64-character, lowercase, hexadecimal strings.
	// Used for validating event IDs and cryptographic keys.
	Hex64Pattern = regexp.MustCompile("^[a-f0-9]{64}$")

	// Hex128Pattern matches 128-character, lowercase, hexadecimal strings.
	// Used for validating signatures.
	Hex128Pattern = regexp.MustCompile("^[a-f0-9]{128}$")
)

var (
	// ErrMalformedPubKey indicates a public key is not 64 lowercase hex characters.
	ErrMalformedPubKey = errors.New("public key must be 64 lowercase hex characters")

	// ErrMalformedPrivKey indicates a private key is not 64 lowercase hex characters.
	ErrMalformedPrivKey = errors.New("private key must be 64 lowercase hex characters")

	// ErrMalformedID indicates an event id is not 64 hex characters.
	ErrMalformedID = errors.New("event id must be 64 hex characters")

	// ErrMalformedSig indicates an event signature is not 128 hex characters.
	ErrMalformedSig = errors.New("event signature must be 128 hex characters")

	// ErrMalformedTag indicates an event tag contains fewer than two elements.
	ErrMalformedTag = errors.New("tags must contain at least two elements")

	// ErrFailedIDComp indicates the event ID could not be computed during validation.
	ErrFailedIDComp = errors.New("failed to compute event id")

	// ErrNoEventID indicates the event ID field is empty.
	ErrNoEventID = errors.New("event id is empty")

	// ErrInvalidSig indicates the event signature failed cryptographic validation.
	ErrInvalidSig = errors.New("event signature is invalid")
)
