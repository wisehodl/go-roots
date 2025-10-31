package events

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// Serialize returns the canonical JSON array representation of the event.
// used for ID computation: [0, pubkey, created_at, kind, tags, content].
func Serialize(e Event) ([]byte, error) {
	serialized := []interface{}{
		0,
		e.PubKey,
		e.CreatedAt,
		e.Kind,
		e.Tags,
		e.Content,
	}

	bytes, err := json.Marshal(serialized)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

// GetID computes and returns the event ID as a lowercase, hex-encoded SHA-256 hash
// of the serialized event.
func GetID(e Event) (string, error) {
	bytes, err := Serialize(e)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:]), nil
}
