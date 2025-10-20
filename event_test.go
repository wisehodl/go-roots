package roots

import (
	"testing"
)

const testSK = "f43a0435f69529f310bbd1d6263d2fbf0977f54bfe2310cc37ae5904b83bb167"
const testPK = "cfa87f35acbde29ba1ab3ee42de527b2cad33ac487e80cf2d6405ea0042c8fef"

var testEvent = Event{
	ID:        "c7a702e6158744ca03508bbb4c90f9dbb0d6e88fefbfaa511d5ab24b4e3c48ad",
	PubKey:    testPK,
	CreatedAt: 1760740551,
	Kind:      1,
	Tags:      [][]string{},
	Content:   "hello world",
	Sig:       "83b71e15649c9e9da362c175f988c36404cabf357a976d869102a74451cfb8af486f6088b5631033b4927bd46cad7a0d90d7f624aefc0ac260364aa65c36071a",
}

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
