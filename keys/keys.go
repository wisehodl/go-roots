package keys

import (
	"encoding/hex"
	"git.wisehodl.dev/jay/go-roots/errors"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// GeneratePrivateKey generates a new, random secp256k1 private key and returns
// it as a 64-character, lowercase hexadecimal string.
func GeneratePrivateKey() (string, error) {
	sk, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return "", err
	}
	skBytes := sk.Serialize()
	return hex.EncodeToString(skBytes), nil
}

// GetPublicKey derives the public key from a private key hex string
// and returns the x-coordinate as 64 lowercase hex characters.
func GetPublicKey(privateKeyHex string) (string, error) {
	if len(privateKeyHex) != 64 {
		return "", errors.MalformedPrivKey
	}
	skBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", errors.MalformedPrivKey
	}

	pk := secp256k1.PrivKeyFromBytes(skBytes).PubKey()
	pkBytes := pk.SerializeCompressed()[1:]
	return hex.EncodeToString(pkBytes), nil
}
