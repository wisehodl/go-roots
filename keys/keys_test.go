package keys

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

const testSK = "f43a0435f69529f310bbd1d6263d2fbf0977f54bfe2310cc37ae5904b83bb167"
const testPK = "cfa87f35acbde29ba1ab3ee42de527b2cad33ac487e80cf2d6405ea0042c8fef"

var Hex64Pattern = regexp.MustCompile("^[a-f0-9]{64}$")

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
