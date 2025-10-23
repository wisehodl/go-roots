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
	assert.Equal(t, actualSig, expectedSig)
}

func TestSignInvalidEventID(t *testing.T) {
	badEventID := "thisisabadeventid"
	expectedError := "invalid event id hex"

	_, err := SignEvent(badEventID, testSK)

	assert.ErrorContains(t, err, expectedError)
}

func TestSignInvalidPrivateKey(t *testing.T) {
	eventID := testEvent.ID
	badSK := "thisisabadsecretkey"
	expectedError := "invalid private key hex"

	_, err := SignEvent(eventID, badSK)

	assert.ErrorContains(t, err, expectedError)
}
