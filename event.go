package roots

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

type Event struct {
	ID        string     `json:"id"`
	PubKey    string     `json:"pubkey"`
	CreatedAt int        `json:"created_at"`
	Kind      int        `json:"kind"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	Sig       string     `json:"sig"`
}

var (
	Hex64Pattern  = regexp.MustCompile("^[a-f0-9]{64}$")
	Hex128Pattern = regexp.MustCompile("^[a-f0-9]{128}$")
)

func (e Event) Serialize() ([]byte, error) {
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

func (e Event) GetID() (string, error) {
	bytes, err := e.Serialize()
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:]), nil
}

func SignEvent(eventID, privateKeyHex string) (string, error) {
	skBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("invalid private key hex: %w", err)
	}

	idBytes, err := hex.DecodeString(eventID)
	if err != nil {
		return "", fmt.Errorf("invalid event id hex: %w", err)
	}

	// discard public key return value
	sk, _ := btcec.PrivKeyFromBytes(skBytes)
	sig, err := schnorr.Sign(sk, idBytes)
	if err != nil {
		return "", fmt.Errorf("schnorr signature error: %w", err)
	}

	return hex.EncodeToString(sig.Serialize()), nil
}

func (e Event) Validate() error {
	if err := e.ValidateStructure(); err != nil {
		return err
	}

	if err := e.ValidateID(); err != nil {
		return err
	}

	return ValidateSignature(e.ID, e.Sig, e.PubKey)
}

func (e Event) ValidateStructure() error {
	if !Hex64Pattern.MatchString(e.PubKey) {
		return fmt.Errorf("pubkey must be 64 lowercase hex characters")
	}

	if !Hex64Pattern.MatchString(e.ID) {
		return fmt.Errorf("id must be 64 hex characters")
	}

	if !Hex128Pattern.MatchString(e.Sig) {
		return fmt.Errorf("signature must be 128 hex characters")
	}

	for _, tag := range e.Tags {
		if len(tag) < 2 {
			return fmt.Errorf("tags must contain at least two elements")
		}
	}

	return nil
}

func (e Event) ValidateID() error {
	computedID, err := e.GetID()
	if err != nil {
		return fmt.Errorf("failed to compute event id")
	}
	if e.ID == "" {
		return fmt.Errorf("event id is empty")
	}
	if computedID != e.ID {
		return fmt.Errorf("event id %q does not match computed id %q", e.ID, computedID)
	}
	return nil
}

func ValidateSignature(eventID, eventSig, publicKeyHex string) error {
	idBytes, err := hex.DecodeString(eventID)
	if err != nil {
		return fmt.Errorf("invalid event id hex: %w", err)
	}

	sigBytes, err := hex.DecodeString(eventSig)
	if err != nil {
		return fmt.Errorf("invalid event signature hex: %w", err)
	}

	pkBytes, err := hex.DecodeString(publicKeyHex)
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
	} else {
		return fmt.Errorf("event signature is invalid")
	}
}
