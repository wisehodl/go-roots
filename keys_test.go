package roots

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var hexPattern = regexp.MustCompile("^[a-f0-9]{64}$")

func TestGeneratePrivateKey(t *testing.T) {
	sk, err := GeneratePrivateKey()

	assert.NoError(t, err)
	if !hexPattern.MatchString(sk) {
		t.Errorf("invalid private key format: %s", sk)
	}
}

func TestGetPublicKey(t *testing.T) {
	pk, err := GetPublicKey(testSK)

	assert.NoError(t, err)
	assert.Equal(t, testPK, pk)
}

func TestGetPublicKeyInvalidPrivateKey(t *testing.T) {
	_, err := GetPublicKey("abc123")
	assert.ErrorContains(t, err, "private key must be 64 lowercase hex characters")
}
