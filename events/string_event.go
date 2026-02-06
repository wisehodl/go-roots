package events

// StringEvent represents a Nostr protocol event, with its seven required
// fields, using hex strings. All fields must be present for a valid event.
type StringEvent struct {
	ID        string `json:"id"`
	PubKey    string `json:"pubkey"`
	CreatedAt int    `json:"created_at"`
	Kind      int    `json:"kind"`
	Tags      []Tag  `json:"tags"`
	Content   string `json:"content"`
	Sig       string `json:"sig"`
}

func (e *StringEvent) GetID() string      { return e.ID }
func (e *StringEvent) GetPubKey() string  { return e.PubKey }
func (e *StringEvent) GetCreatedAt() int  { return e.CreatedAt }
func (e *StringEvent) GetKind() int       { return e.Kind }
func (e *StringEvent) GetTags() []Tag     { return e.Tags }
func (e *StringEvent) GetContent() string { return e.Content }
func (e *StringEvent) GetSig() string     { return e.Sig }

func (e *StringEvent) SetID(id string)   { e.ID = id }
func (e *StringEvent) SetSig(sig string) { e.Sig = sig }
