package roots

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

// SignEvent generates a Schnorr signature for the given event ID using the
// provided private key. Returns the signature as 128 lowercase hex characters.
func SignEvent(eventID, privateKeyHex string) (string, error) {
	skBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", ErrMalformedPrivKey
	}

	idBytes, err := hex.DecodeString(eventID)
	if err != nil {
		return "", ErrMalformedID
	}

	// discard public key return value
	sk, _ := btcec.PrivKeyFromBytes(skBytes)
	sig, err := schnorr.Sign(sk, idBytes)
	if err != nil {
		return "", fmt.Errorf("schnorr signature error: %w", err)
	}

	return hex.EncodeToString(sig.Serialize()), nil
}
