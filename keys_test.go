package roots

import (
	"regexp"
	"testing"
)

var hexPattern = regexp.MustCompile("^[a-f0-9]{64}$")

func TestGeneratePrivateKey(t *testing.T) {
	sk, err := GeneratePrivateKey()

	expectOk(t, err)
	if !hexPattern.MatchString(sk) {
		t.Errorf("invalid private key format: %s", sk)
	}
}

func TestGetPublicKey(t *testing.T) {
	pk, err := GetPublicKey(testSK)

	expectOk(t, err)
	expectEqualStrings(t, pk, testPK)
}

func TestGetPublicKeyInvalidPrivateKey(t *testing.T) {
	_, err := GetPublicKey("abc123")
	expectErrorSubstring(t, err, "private key must be 64 hex characters")
}
