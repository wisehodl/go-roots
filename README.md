# go-roots — Nostr Protocol Library for Go

Source: https://git.wisehodl.dev/jay/go-roots

Mirror: https://github.com/wisehodl/go-roots

## What this library does

`go-roots` is a consensus-layer Nostr protocol library for Go.
It provides primitives that define protocol compliance:

- Event structure and serialization
- Cryptographic signing and validation
- Subscription filters

## What this library does not do

`go-roots` serves as a foundation for libraries and applications that implement higher-level abstractions of the Nostr protocol, including message transport, semantic event definitions, event storage, and user interfaces.

`go-roots` prioritizes correctness and clarity over optimization and efficiency. High-performance applications should implement optimizations in a separate library or in the application itself.

## Installation

1. Add `go-roots` to your project:

```bash
go get git.wisehodl.dev/jay/go-roots
```

If the primary repository is unavailable, use the `replace` directive in your go.mod file to get the package from the GitHub mirror:

```
replace git.wisehodl.dev/jay/go-roots => github.com/wisehodl/go-roots latest
```

2. Import the packages:

```go
import (
    "git.wisehodl.dev/jay/go-roots/errors"
    "git.wisehodl.dev/jay/go-roots/events"
    "git.wisehodl.dev/jay/go-roots/filters"
    "git.wisehodl.dev/jay/go-roots/keys"
)
```

## General Use

`ValidatedEvent` is the primary type consumers work with. A plain `Event` is a mutable scratch pad used during construction or deserialization. Once an event passes cryptographic validation, it is promoted to `ValidatedEvent` through a single gate: `NewValidatedEvent`. After that point, the event is immutable and its fields are accessed through methods.

### Key Management

#### Generate a new keypair

```go
privateKey, err := keys.GeneratePrivateKey()
if err != nil {
    log.Fatal(err)
}

publicKey, err := keys.GetPublicKey(privateKey)
if err != nil {
    log.Fatal(err)
}
```

#### Derive public key from an existing private key

```go
privateKey := "f43a0435f69529f310bbd1d6263d2fbf0977f54bfe2310cc37ae5904b83bb167"
publicKey, err := keys.GetPublicKey(privateKey)
// publicKey: "cfa87f35acbde29ba1ab3ee42de527b2cad33ac487e80cf2d6405ea0042c8fef"
```

### Flow 1: Creating an event

Build the event with `NewEvent`, compute the ID with `GetID`, sign it with `SignEvent`, then promote to `ValidatedEvent` with `NewValidatedEvent`.

```go
// 1. Build the event
event := events.NewEvent(
    events.WithPubKey(publicKey),
    events.WithCreatedAt(time.Now().Unix()),
    events.WithKind(1),
    events.WithTag(events.Tag{"e", "5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36"}),
    events.WithTag(events.Tag{"p", "91cf9b32f3735070f46c0a86a820a47efa08a5be6c9f4f8cf68e5b5b75c92d60"}),
    events.WithContent("Hello, Nostr!"),
)

// 2. Compute and assign the ID
event.ID = events.GetID(event)

// 3. Sign the event
sig, err := events.SignEvent(event.ID, privateKey)
if err != nil {
    log.Fatal(err)
}
event.Sig = sig

// 4. Promote to ValidatedEvent
validated, err := events.NewValidatedEvent(event)
if err != nil {
    log.Fatal(err) // construction bug — treat as a defect, not a runtime condition
}
```

For applications that delegate signing to a remote signer (NIP-46, hardware device, custody service): build the unsigned `Event`, send it to the signer, receive the signed `Event` back, then call `NewValidatedEvent` on what you received.

### Flow 2: Receiving an event from an external source

Unmarshal into a plain `Event`, then immediately promote to `ValidatedEvent` with `NewValidatedEvent`.

```go
var event events.Event
if err := json.Unmarshal(data, &event); err != nil {
    log.Printf("Malformed JSON: %v", err)
    return
}

validated, err := events.NewValidatedEvent(event)
if err != nil {
    log.Printf("Invalid event: %v", err) // reject from peer; investigate if from database
    return
}
```

### Flow 3: Serializing a validated event to JSON

Call `json.Marshal` on a `ValidatedEvent`. The output is guaranteed to reflect a cryptographically verified event.

```go
jsonBytes, err := json.Marshal(validated)
if err != nil {
    log.Fatal(err)
}
```

Calling `json.Marshal` on a plain `Event` is discouraged: the resulting JSON carries no integrity guarantee.

### The `NewValidatedEvent` gate

`NewValidatedEvent` is the single promotion point. It verifies:

- Valid field structure (hex lengths, tag shape, field formats)
- ID matches the SHA-256 hash of the canonical serialization
- Signature is a valid Schnorr signature over the ID for the given public key

It does not verify semantic validity. A `ValidatedEvent` is valid according to the base Nostr protocol. Whether it is valid for your application — correct kind, trusted author, within a relevant time window — is your responsibility.

### Filters

#### Build a filter

```go
since := int(time.Now().Add(-24 * time.Hour).Unix())
limit := 50

filter := filters.Filter{
    IDs:     []string{"abc123", "def456"},  // prefix match
    Authors: []string{"cfa87f35"},          // prefix match
    Kinds:   []int{1, 6, 7},
    Since:   &since,
    Limit:   &limit,
}
```

#### Tag filters

```go
filter := filters.Filter{
    Kinds: []int{1},
    Tags: filters.TagFilters{
        "e": {"5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36"},
        "p": {"91cf9b32f3735070f46c0a86a820a47efa08a5be6c9f4f8cf68e5b5b75c92d60"},
    },
}
```

#### Match events against a filter

`Matches` accepts a `ValidatedEvent`. Unvalidated events cannot be passed to it.

```go
filter := filters.Filter{
    Authors: []string{"cfa87f35"},
    Kinds:   []int{1},
}

var matched []events.ValidatedEvent
for _, event := range incoming {
    if filters.Matches(filter, event) {
        matched = append(matched, event)
    }
}
```

`Matches` covers the standard NIP-01 filter fields. For non-trivial matching logic, implement your own.

#### Filter JSON

Tag filter keys are prefixed with `#` in JSON (`#e`, `#p`). Unrecognized fields are captured in `Extensions` and are retained as-found. `Matches` ignores `Extensions`; higher layers can inspect them for custom behavior.

```go
// Marshal
jsonBytes, err := filters.MarshalJSON(filter)
// {"authors":["cfa87f35"],"kinds":[1],"#e":["..."],"since":...}

// Unmarshal
var f filters.Filter
if err := filters.UnmarshalJSON(data, &f); err != nil {
    log.Fatal(err)
}
```

Extensions example:

```go
filter := filters.Filter{
    Kinds: []int{1},
    Extensions: filters.FilterExtensions{
        "search": json.RawMessage(`"bitcoin"`),
    },
}

// In a storage layer (not this library):
if searchRaw, ok := filter.Extensions["search"]; ok {
    var term string
    json.Unmarshal(searchRaw, &term)
    // apply full-text search
}
```

---

## Specialized / Low-Level Tools

These are building blocks for custom pipelines. General use does not require them.

- **`Validate(e Event) error`** — full validation in one call; used internally by `NewValidatedEvent`
- **`ValidateStructure(e Event) error`** — field format checks only
- **`ValidateSignature(e Event) error`** — Schnorr signature verification only
- **`Serialize(e Event) []byte`** — canonical `[0, pubkey, created_at, kind, tags, content]` JSON array; used for ID computation
- **`GetID(e Event) string`** — SHA-256 of `Serialize`, returned as lowercase hex
- **`IsValidKey(s string) bool`**, **`IsValidID`**, **`IsValidSig`** — format validators for individual field values; useful when handling these value types outside of event validation

## Testing

```bash
go test ./...
```
