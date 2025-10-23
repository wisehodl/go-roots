package roots

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignEvent(t *testing.T) {
	eventID := testEvent.ID
	expectedSig := testEvent.Sig
	actualSig, err := SignEvent(eventID, testSK)

	assert.NoError(t, err)
	assert.Equal(t, expectedSig, actualSig)
}

func TestSignInvalidEventID(t *testing.T) {
	badEventID := "thisisabadeventid"
	expectedError := "event id must be 64 hex characters"

	_, err := SignEvent(badEventID, testSK)

	assert.ErrorContains(t, err, expectedError)
}

func TestSignInvalidPrivateKey(t *testing.T) {
	eventID := testEvent.ID
	badSK := "thisisabadsecretkey"
	expectedError := "private key must be 64 lowercase hex characters"

	_, err := SignEvent(eventID, badSK)

	assert.ErrorContains(t, err, expectedError)
}
