package roots

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	sk, err := GeneratePrivateKey()

	assert.NoError(t, err)
	if !Hex64Pattern.MatchString(sk) {
		t.Errorf("invalid private key format: %s", sk)
	}
}

func TestGenerateUniquePrivateKeys(t *testing.T) {
	sk1, _ := GeneratePrivateKey()
	sk2, _ := GeneratePrivateKey()
	assert.NotEqual(t, sk1, sk2)
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
