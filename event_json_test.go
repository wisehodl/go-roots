package roots

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalEventJSON(t *testing.T) {
	event := Event{}
	json.Unmarshal(testEventJSONBytes, &event)
	if err := event.Validate(); err != nil {
		t.Error("unmarshalled event is invalid")
	}
	expectEqualEvents(t, event, testEvent)
}

func TestMarshalEventJSON(t *testing.T) {
	eventJSONBytes, err := json.Marshal(testEvent)
	expectOk(t, err)
	expectEqualStrings(t, string(eventJSONBytes), testEventJSON)
}

func TestEventJSONRoundTrip(t *testing.T) {
	event := Event{
		ID:        "86e856d0527dd08527498cd8afd8a7d296bde37e4757a8921f034f0b344df3ad",
		PubKey:    testEvent.PubKey,
		CreatedAt: testEvent.CreatedAt,
		Kind:      testEvent.Kind,
		Tags: [][]string{
			{"a", "value"},
			{"b", "value", "optional"},
			{"name", "value", "optional", "optional"},
		},
		Content: testEvent.Content,
		Sig:     "c05fe02a9c082ff56aad2b16b5347498a21665f02f050ba086dbe6bd593c8cd448505d2831d1c0340acc1793eaf89b7c0cb21bb696c71da6b8d6b857702bb557",
	}
	expectedJSON := `{"id":"86e856d0527dd08527498cd8afd8a7d296bde37e4757a8921f034f0b344df3ad","pubkey":"cfa87f35acbde29ba1ab3ee42de527b2cad33ac487e80cf2d6405ea0042c8fef","created_at":1760740551,"kind":1,"tags":[["a","value"],["b","value","optional"],["name","value","optional","optional"]],"content":"hello world","sig":"c05fe02a9c082ff56aad2b16b5347498a21665f02f050ba086dbe6bd593c8cd448505d2831d1c0340acc1793eaf89b7c0cb21bb696c71da6b8d6b857702bb557"}`

	if err := event.Validate(); err != nil {
		t.Error("test event is invalid")
	}
	eventJSON, err := json.Marshal(event)
	expectOk(t, err)
	expectEqualStrings(t, string(eventJSON), expectedJSON)

	unmarshalledEvent := Event{}
	json.Unmarshal(eventJSON, &unmarshalledEvent)

	if err := unmarshalledEvent.Validate(); err != nil {
		t.Error("unmarshalled event is invalid")
	}
	expectEqualEvents(t, unmarshalledEvent, event)
}
