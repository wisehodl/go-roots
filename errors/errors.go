package errors

import (
	"errors"
)

var (
	// MalformedPubKey indicates a public key is not 64 lowercase hex characters.
	MalformedPubKey = errors.New("public key must be 64 lowercase hex characters")

	// MalformedPrivKey indicates a private key is not 64 lowercase hex characters.
	MalformedPrivKey = errors.New("private key must be 64 lowercase hex characters")

	// MalformedID indicates an event id is not 64 hex characters.
	MalformedID = errors.New("event id must be 64 hex characters")

	// MalformedSig indicates an event signature is not 128 hex characters.
	MalformedSig = errors.New("event signature must be 128 hex characters")

	// MalformedTag indicates an event tag contains fewer than two elements.
	MalformedTag = errors.New("tags must contain at least two elements")

	// FailedIDComp indicates the event ID could not be computed during validation.
	FailedIDComp = errors.New("failed to compute event id")

	// NoEventID indicates the event ID field is empty.
	NoEventID = errors.New("event id is empty")

	// InvalidSig indicates the event signature failed cryptographic validation.
	InvalidSig = errors.New("event signature is invalid")
)
