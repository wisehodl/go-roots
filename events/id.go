package events

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

// Serialize returns the canonical JSON array representation of the event.
// used for ID computation: [0, pubkey, created_at, kind, tags, content].
func Serialize(e Event) []byte {
	buf := make([]byte, 0, 100+len(e.Content)+len(e.Tags)*80)
	return appendSerialized(buf, e)
}

// GetID computes and returns the event ID as a lowercase, hex-encoded SHA-256 hash
// of the serialized event.
func GetID(e Event) string {
	bytes := Serialize(e)
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:])
}

func appendSerialized(dst []byte, e Event) []byte {
	dst = append(dst, "[0,\""...)
	dst = append(dst, e.PubKey...)
	dst = append(dst, "\","...)
	dst = strconv.AppendInt(dst, int64(e.CreatedAt), 10)
	dst = append(dst, ',')
	dst = strconv.AppendInt(dst, int64(e.Kind), 10)
	dst = append(dst, ",["...)
	for i, tag := range e.Tags {
		if i > 0 {
			dst = append(dst, ',')
		}
		dst = append(dst, '[')
		for j, s := range tag {
			if j > 0 {
				dst = append(dst, ',')
			}
			dst = appendEscapedString(dst, s)
		}
		dst = append(dst, ']')
	}
	dst = append(dst, "],"...)
	dst = appendEscapedString(dst, e.Content)
	return append(dst, ']')
}

func appendEscapedString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			dst = append(dst, '\\', '"')
		case c == '\\':
			dst = append(dst, '\\', '\\')
		case c >= 0x20:
			dst = append(dst, c)
		case c == 0x08:
			dst = append(dst, '\\', 'b')
		case c == 0x09:
			dst = append(dst, '\\', 't')
		case c == 0x0a:
			dst = append(dst, '\\', 'n')
		case c == 0x0c:
			dst = append(dst, '\\', 'f')
		case c == 0x0d:
			dst = append(dst, '\\', 'r')
		case c < 0x09:
			dst = append(dst, '\\', 'u', '0', '0', '0', '0'+c)
		case c < 0x10:
			dst = append(dst, '\\', 'u', '0', '0', '0', 0x57+c)
		case c < 0x1a:
			dst = append(dst, '\\', 'u', '0', '0', '1', 0x20+c)
		case c < 0x20:
			dst = append(dst, '\\', 'u', '0', '0', '1', 0x47+c)
		}
	}
	return append(dst, '"')
}
