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
		event: Event{
			ID:        testEvent.ID,
			PubKey:    "",
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      testEvent.Tags,
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "short pubkey",
		event: Event{
			ID:        testEvent.ID,
			PubKey:    "abc123",
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      testEvent.Tags,
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "long pubkey",
		event: Event{
			ID:        testEvent.ID,
			PubKey:    "c7a702e6158744ca03508bbb4c90f9dbb0d6e88fefbfaa511d5ab24b4e3c48adabc",
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      testEvent.Tags,
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "non-hex pubkey",
		event: Event{
			ID:        testEvent.ID,
			PubKey:    "zyx-!2e6158744ca03508bbb4c90f9dbb0d6e88fefbfaa511d5ab24b4e3c48ad",
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      testEvent.Tags,
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "uppercase pubkey",
		event: Event{
			ID:        testEvent.ID,
			PubKey:    "C7A702E6158744CA03508BBB4C90F9DBB0D6E88FEFBFAA511D5AB24B4E3C48AD",
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      testEvent.Tags,
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "public key must be 64 lowercase hex characters",
	},

	{
		name: "empty id",
		event: Event{
			ID:        "",
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      testEvent.Tags,
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "id must be 64 hex characters",
	},

	{
		name: "short id",
		event: Event{
			ID:        "abc123",
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      testEvent.Tags,
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "id must be 64 hex characters",
	},

	{
		name: "empty signature",
		event: Event{
			ID:        testEvent.ID,
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      testEvent.Tags,
			Content:   testEvent.Content,
			Sig:       "",
		},
		expectedError: "signature must be 128 hex characters",
	},

	{
		name: "short signature",
		event: Event{
			ID:        testEvent.ID,
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      testEvent.Tags,
			Content:   testEvent.Content,
			Sig:       "abc123",
		},
		expectedError: "signature must be 128 hex characters",
	},

	{
		name: "empty tag",
		event: Event{
			ID:        testEvent.ID,
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      []Tag{{}},
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "tags must contain at least two elements",
	},

	{
		name: "single element tag",
		event: Event{
			ID:        testEvent.ID,
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      []Tag{{"a"}},
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "tags must contain at least two elements",
	},

	{
		name: "one good tag, one single element tag",
		event: Event{
			ID:        testEvent.ID,
			PubKey:    testEvent.PubKey,
			CreatedAt: testEvent.CreatedAt,
			Kind:      testEvent.Kind,
			Tags:      []Tag{{"a", "value"}, {"b"}},
			Content:   testEvent.Content,
			Sig:       testEvent.Sig,
		},
		expectedError: "tags must contain at least two elements",
	},
}

func TestValidateEventStructure(t *testing.T) {
	for _, tc := range structureTestCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.event.ValidateStructure()
			assert.ErrorContains(t, err, tc.expectedError)
		})
	}
}

func TestValidateEventIDFailure(t *testing.T) {
	event := Event{
		ID:        "7f661c2a3c1ed67dc959d6cd968d743d5e6e334313df44724bca939e2aa42c9e",
		PubKey:    testEvent.PubKey,
		CreatedAt: testEvent.CreatedAt,
		Kind:      testEvent.Kind,
		Tags:      testEvent.Tags,
		Content:   testEvent.Content,
		Sig:       testEvent.Sig,
	}

	err := event.ValidateID()
	assert.ErrorContains(t, err, "does not match computed id")
}

func TestValidateSignature(t *testing.T) {
	event := Event{
		ID:     testEvent.ID,
		PubKey: testEvent.PubKey,
		Sig:    testEvent.Sig,
	}
	err := event.ValidateSignature()

	assert.NoError(t, err)
}

func TestValidateInvalidSignature(t *testing.T) {
	event := Event{
		ID:     testEvent.ID,
		PubKey: testEvent.PubKey,
		Sig:    "9e43cbcf7e828a21c53fa35371ee79bffbfd7a3063ae46fc05ec623dd3186667c57e3d006488015e19247df35eb41c61013e051aa87860e23fa5ffbd44120482",
	}
	err := event.ValidateSignature()

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
			event := Event{ID: tc.id, PubKey: tc.pubkey, Sig: tc.sig}
			err := event.ValidateSignature()
			assert.ErrorContains(t, err, tc.expectedError)
		})
	}
}

func TestValidateEvent(t *testing.T) {
	event := Event{
		ID:        "c9a0f84fcaa889654da8992105eb122eb210c8cbd58210609a5ef7e170b51400",
		PubKey:    testEvent.PubKey,
		CreatedAt: testEvent.CreatedAt,
		Kind:      testEvent.Kind,
		Tags: []Tag{
			{"a", "value"},
			{"b", "value", "optional"},
		},
		Content: "valid event",
		Sig:     "668a715f1eb983172acf230d17bd283daedb2598adf8de4290bcc7eb0b802fdb60669d1e7d1104ac70393f4dbccd07e8abf897152af6ce6c0a75499874e27f14",
	}

	err := event.Validate()
	assert.NoError(t, err)
}
