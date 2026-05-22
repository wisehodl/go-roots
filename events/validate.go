package events

import (
	"bytes"
	"crypto/sha256"
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

	idHash := sha256.Sum256(Serialize(e))
	idBytes, err := hex.DecodeString(e.ID)
	if err != nil {
		return errors.MalformedID
	}
	if !bytes.Equal(idBytes, idHash[:]) {
		return fmt.Errorf(
			"event id %q does not match computed id %q",
			e.ID, hex.EncodeToString(idHash[:]))
	}

	return validateSignatureBytes(idBytes, e.Sig, e.PubKey)
}

// ValidateStructure checks that all event fields conform to the protocol
// specification: hex lengths, tag structure, and field formats.
func ValidateStructure(e Event) error {
	if !IsValidKey(e.PubKey) {
		return errors.MalformedPubKey
	}

	if !IsValidID(e.ID) {
		return errors.MalformedID
	}

	if !IsValidSig(e.Sig) {
		return errors.MalformedSig
	}

	for _, tag := range e.Tags {
		if len(tag) < 2 {
			return errors.MalformedTag
		}
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
	return validateSignatureBytes(idBytes, e.Sig, e.PubKey)
}

// Value validators

// IsValidKey verifies that a public or private key is properly formatted.
func IsValidKey(value string) bool {
	return isLowerHex(value, 64)
}

// IsValidID verifies that an event id is properly formatted.
func IsValidID(value string) bool {
	return isLowerHex(value, 64)
}

// IsValidSig verifies that an event signature is properly formatted.
func IsValidSig(value string) bool {
	return isLowerHex(value, 128)
}

// Helpers

func validateSignatureBytes(idBytes []byte, sigHex, pkHex string) error {
	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return fmt.Errorf("invalid event signature hex: %w", err)
	}

	pkBytes, err := hex.DecodeString(pkHex)
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
	}

	return errors.InvalidSig
}

func isLowerHex(s string, n int) bool {
	if len(s) != n {
		return false
	}
	for i := 0; i < n; i++ {
		c := s[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}
