package roots

import (
	"encoding/hex"
	"fmt"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func GeneratePrivateKey() (string, error) {
	sk, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return "", err
	}
	skBytes := sk.Serialize()
	return hex.EncodeToString(skBytes), nil
}

func GetPublicKey(privateKeyHex string) (string, error) {
	if len(privateKeyHex) != 64 {
		return "", fmt.Errorf("private key must be 64 hex characters")
	}
	skBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}

	pk := secp256k1.PrivKeyFromBytes(skBytes).PubKey()
	pkBytes := pk.SerializeCompressed()[1:]
	return hex.EncodeToString(pkBytes), nil
}
