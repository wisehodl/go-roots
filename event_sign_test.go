package roots

import (
	"testing"
)

func TestSignEvent(t *testing.T) {
	eventID := testEvent.ID
	expectedSig := testEvent.Sig
	computedSig, err := SignEvent(eventID, testSK)

	expectOk(t, err)
	expectEqualStrings(t, computedSig, expectedSig)
}

func TestSignInvalidEventID(t *testing.T) {
	badEventID := "thisisabadeventid"
	expectedError := "invalid event id hex"

	_, err := SignEvent(badEventID, testSK)

	expectError(t, err)
	expectErrorSubstring(t, err, expectedError)
}

func TestSignInvalidPrivateKey(t *testing.T) {
	eventID := testEvent.ID
	badSK := "thisisabadsecretkey"
	expectedError := "invalid private key hex"

	_, err := SignEvent(eventID, badSK)

	expectError(t, err)
	expectErrorSubstring(t, err, expectedError)
}
