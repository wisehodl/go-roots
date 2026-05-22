package events

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type ValidateEventTestCase struct {
	name          string
	event         Event
	expectedError string
}

var structureTestCases = []ValidateEventTestCase{
	{
		name: "empty pubkey",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey(""),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "short pubkey",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey("abc123"),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "long pubkey",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey("c7a702e6158744ca03508bbb4c90f9dbb0d6e88fefbfaa511d5ab24b4e3c48adabc"),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "non-hex pubkey",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey("zyx-!2e6158744ca03508bbb4c90f9dbb0d6e88fefbfaa511d5ab24b4e3c48ad"),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "uppercase pubkey",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey("C7A702E6158744CA03508BBB4C90F9DBB0D6E88FEFBFAA511D5AB24B4E3C48AD"),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "empty id",
		event: NewEvent(
			WithID(""),
			WithPubKey(testEvent.PubKey),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "id must be 64 hex characters",
	},

	{
		name: "short id",
		event: NewEvent(
			WithID("abc123"),
			WithPubKey(testEvent.PubKey),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "id must be 64 hex characters",
	},

	{
		name: "empty signature",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey(testEvent.PubKey),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig(""),
		),
		expectedError: "signature must be 128 hex characters",
	},

	{
		name: "short signature",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey(testEvent.PubKey),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithContent(testEvent.Content),
			WithSig("abc123"),
		),
		expectedError: "signature must be 128 hex characters",
	},

	{
		name: "empty tag",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey(testEvent.PubKey),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithTag(Tag{}),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "tags must contain at least two elements",
	},

	{
		name: "single element tag",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey(testEvent.PubKey),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithTag(Tag{"a"}),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "tags must contain at least two elements",
	},

	{
		name: "one good tag, one single element tag",
		event: NewEvent(
			WithID(testEvent.ID),
			WithPubKey(testEvent.PubKey),
			WithCreatedAt(testEvent.CreatedAt),
			WithKind(testEvent.Kind),
			WithTag(Tag{"a", "value"}),
			WithTag(Tag{"b"}),
			WithContent(testEvent.Content),
			WithSig(testEvent.Sig),
		),
		expectedError: "tags must contain at least two elements",
	},
}

func TestValidateEventStructure(t *testing.T) {
	for _, tc := range structureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateStructure(tc.event)
			assert.ErrorContains(t, err, tc.expectedError)
		})
	}
}

func TestValidateSignature(t *testing.T) {
	event := NewEvent(
		WithID(testEvent.ID),
		WithPubKey(testEvent.PubKey),
		WithSig(testEvent.Sig),
	)
	err := ValidateSignature(event)

	assert.NoError(t, err)
}

func TestValidateInvalidSignature(t *testing.T) {
	event := NewEvent(
		WithID(testEvent.ID),
		WithPubKey(testEvent.PubKey),
		WithSig("9e43cbcf7e828a21c53fa35371ee79bffbfd7a3063ae46fc05ec623dd3186667c57e3d006488015e19247df35eb41c61013e051aa87860e23fa5ffbd44120482"),
	)
	err := ValidateSignature(event)

	assert.ErrorContains(t, err, "event signature is invalid")
}

type ValidateSignatureTestCase struct {
	name          string
	id            string
	sig           string
	pubkey        string
	expectedError string
}

var validateSignatureTestCases = []ValidateSignatureTestCase{
	{
		name:          "bad event id",
		id:            "badeventid",
		sig:           testEvent.Sig,
		pubkey:        testEvent.PubKey,
		expectedError: "invalid event id hex",
	},

	{
		name:          "bad event signature",
		id:            testEvent.ID,
		sig:           "badeventsignature",
		pubkey:        testEvent.PubKey,
		expectedError: "invalid event signature hex",
	},

	{
		name:          "bad public key",
		id:            testEvent.ID,
		sig:           testEvent.Sig,
		pubkey:        "badpublickey",
		expectedError: "invalid public key hex",
	},

	{
		name:          "malformed event signature",
		id:            testEvent.ID,
		sig:           "abc123",
		pubkey:        testEvent.PubKey,
		expectedError: "malformed signature",
	},

	{
		name:          "malformed public key",
		id:            testEvent.ID,
		sig:           testEvent.Sig,
		pubkey:        "abc123",
		expectedError: "malformed public key",
	},
}

func TestValidateSignatureInvalidEventSignature(t *testing.T) {
	for _, tc := range validateSignatureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			event := NewEvent(
				WithID(tc.id),
				WithPubKey(tc.pubkey),
				WithSig(tc.sig),
			)
			err := ValidateSignature(event)
			assert.ErrorContains(t, err, tc.expectedError)
		})
	}
}

func TestValidateEvent(t *testing.T) {
	event := NewEvent(
		WithID("c9a0f84fcaa889654da8992105eb122eb210c8cbd58210609a5ef7e170b51400"),
		WithPubKey(testEvent.PubKey),
		WithCreatedAt(testEvent.CreatedAt),
		WithKind(testEvent.Kind),
		WithTag(Tag{"a", "value"}),
		WithTag(Tag{"b", "value", "optional"}),
		WithContent("valid event"),
		WithSig("668a715f1eb983172acf230d17bd283daedb2598adf8de4290bcc7eb0b802fdb60669d1e7d1104ac70393f4dbccd07e8abf897152af6ce6c0a75499874e27f14"),
	)

	err := Validate(event)
	assert.NoError(t, err)
}
