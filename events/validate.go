package events

import (
	"encoding/hex"
	"fmt"
	"git.wisehodl.dev/jay/go-roots/errors"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

// Validate performs a complete event validation: structure, ID computation,
// and signature verification. Returns the first error encountered.
func Validate(e Event) error {
	if err := ValidateStructure(e); err != nil {
		return err
	}

	if err := ValidateID(e); err != nil {
		return err
	}

	return ValidateSignature(e)
}

// ValidateStructure checks that all event fields conform to the protocol
// specification: hex lengths, tag structure, and field formats.
func ValidateStructure(e Event) error {
	if !Hex64Pattern.MatchString(e.PubKey) {
		return errors.MalformedPubKey
	}

	if !Hex64Pattern.MatchString(e.ID) {
		return errors.MalformedID
	}

	if !Hex128Pattern.MatchString(e.Sig) {
		return errors.MalformedSig
	}

	for _, tag := range e.Tags {
		if len(tag) < 2 {
			return errors.MalformedTag
		}
	}

	return nil
}

// ValidateID recomputes the event ID and verifies it matches the stored ID field.
func ValidateID(e Event) error {
	computedID, err := GetID(e)
	if err != nil {
		return errors.FailedIDComp
	}
	if e.ID == "" {
		return errors.NoEventID
	}
	if computedID != e.ID {
		return fmt.Errorf("event id %q does not match computed id %q", e.ID, computedID)
	}
	return nil
}

// ValidateSignature verifies the event signature is cryptographically valid
// for the event ID and public key using Schnorr verification.
func ValidateSignature(e Event) error {
	idBytes, err := hex.DecodeString(e.ID)
	if err != nil {
		return fmt.Errorf("invalid event id hex: %w", err)
	}

	sigBytes, err := hex.DecodeString(e.Sig)
	if err != nil {
		return fmt.Errorf("invalid event signature hex: %w", err)
	}

	pkBytes, err := hex.DecodeString(e.PubKey)
	if err != nil {
		return fmt.Errorf("invalid public key hex: %w", err)
	}

	signature, err := schnorr.ParseSignature(sigBytes)
	if err != nil {
		return fmt.Errorf("malformed signature: %w", err)
	}

	publicKey, err := schnorr.ParsePubKey(pkBytes)
	if err != nil {
		return fmt.Errorf("malformed public key: %w", err)
	}

	if signature.Verify(idBytes, publicKey) {
		return nil
	} else {
		return errors.InvalidSig
	}
}
