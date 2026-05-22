package events

// ValidatedEvent is an immutable wrapper around a fully validated event. It
// shares no memory with the original Event struct.
//
// When created with NewValidatedEvent, the wrapped event is guaranteed to:
// - have a valid structure
// - have a valid ID
// - have a valid signature
type ValidatedEvent struct {
	event Event
}

// NewValidatedEvent validates the provided Event. If valid, returns an
// immutable ValidatedEvent. Otherwise returns an error.
func NewValidatedEvent(e Event) (ValidatedEvent, error) {
	if err := Validate(e); err != nil {
		return ValidatedEvent{}, err
	}
	e.Tags = deepCopyTags(e.Tags)
	return ValidatedEvent{event: e}, nil
}

func (v ValidatedEvent) ID() string       { return v.event.ID }
func (v ValidatedEvent) PubKey() string   { return v.event.PubKey }
func (v ValidatedEvent) CreatedAt() int64 { return v.event.CreatedAt }
func (v ValidatedEvent) Kind() int        { return v.event.Kind }
func (v ValidatedEvent) Content() string  { return v.event.Content }
func (v ValidatedEvent) Sig() string      { return v.event.Sig }

// Tags returns a deep copy of the event's tag list. The returned slice shares
// no memory with the ValidatedEvent.
func (v ValidatedEvent) Tags() []Tag { return deepCopyTags(v.event.Tags) }

// Event returns a deep copy of the underlying Event. The returned Event shares
// no memory with the ValidatedEvent.
func (v ValidatedEvent) Event() Event {
	e := v.event
	e.Tags = deepCopyTags(v.event.Tags)
	return e
}

func deepCopyTags(tags []Tag) []Tag {
	cp := make([]Tag, len(tags))
	for i, tag := range tags {
		tagcp := make(Tag, len(tag))
		copy(tagcp, tag)
		cp[i] = tagcp
	}
	return cp
}

func (v ValidatedEvent) MarshalJSON() ([]byte, error) {
	return v.event.MarshalJSON()
}
