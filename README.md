# Go-Roots - Nostr Protocol Library for Golang

Source: https://git.wisehodl.dev/jay/go-roots

Mirror: https://github.com/wisehodl/go-roots

## What this library does

`go-roots` is a purposefully minimal Nostr protocol library for golang.
It only provides primitives that define protocol compliance:

- Event Structure
- Serialization
- Cryptographic Signatures
- Subscription Filters

## What this library does not do

`go-roots` serves a foundation for other libraries and applications to
implement higher level abstractions of the Nostr protocol on top of it,
including message transport, semantic event definitions, event storage
mechanisms, and user interfaces.

`go-roots` prioritizes correctness and clarity over optimization and efficiency. For high performance applications, it is recommended to implement optimizations in a separate library or in the application which requires them.

## Installation

1. Add `go-roots` to your project:

```bash
go get git.wisehodl.dev/jay/go-roots
```

If the primary repository is unavailable, use the `replace` directive in your go.mod file to get the package from the github mirror:

```
replace git.wisehodl.dev/jay/go-roots => github.com/wisehodl/go-roots latest
```

2. Import the packages:

```golang
import (
    "git.wisehodl.dev/jay/go-roots/errors"
    "git.wisehodl.dev/jay/go-roots/events"
    "git.wisehodl.dev/jay/go-roots/filters"
    "git.wisehodl.dev/jay/go-roots/keys"
)
```

3. Access functions with appropriate namespaces.

## Usage Examples

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

#### Derive public key from existing private key

```go
privateKey := "f43a0435f69529f310bbd1d6263d2fbf0977f54bfe2310cc37ae5904b83bb167"
publicKey, err := keys.GetPublicKey(privateKey)
// publicKey: "cfa87f35acbde29ba1ab3ee42de527b2cad33ac487e80cf2d6405ea0042c8fef"
```

---

### Event Creation and Signing

#### Create and sign a complete event

```go
// 1. Build the event structure
event := events.Event{
    PubKey:    publicKey,
    CreatedAt: int(time.Now().Unix()),
    Kind:      1,
    Tags: []events.Tag{
        {"e", "5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36"},
        {"p", "91cf9b32f3735070f46c0a86a820a47efa08a5be6c9f4f8cf68e5b5b75c92d60"},
    },
    Content: "Hello, Nostr!",
}

// 2. Compute the event ID
id, err := events.GetID(event)
if err != nil {
    log.Fatal(err)
}
event.ID = id

// 3. Sign the event
sig, err := events.SignEvent(id, privateKey)
if err != nil {
    log.Fatal(err)
}
event.Sig = sig
```

#### Serialize an event for ID computation

```go
// Returns canonical JSON: [0, pubkey, created_at, kind, tags, content]
serialized, err := events.Serialize(event)
if err != nil {
    log.Fatal(err)
}
```

#### Compute event ID manually

```go
id, err := events.GetID(event)
if err != nil {
    log.Fatal(err)
}
// Returns lowercase hex SHA-256 hash of serialized form
```

---

### Event Validation

#### Validate complete event

```go
// Checks structure, ID computation, and signature
if err := events.Validate(event); err != nil {
    log.Printf("Invalid event: %v", err)
}
```

#### Validate individual aspects

```go
// Check field formats and lengths
if err := events.ValidateStructure(event); err != nil {
    log.Printf("Malformed structure: %v", err)
}

// Verify ID matches computed hash
if err := events.ValidateID(event); err != nil {
    log.Printf("ID mismatch: %v", err)
}

// Verify cryptographic signature
if err := events.ValidateSignature(event); err != nil {
    log.Printf("Invalid signature: %v", err)
}
```

---

### Event JSON

#### Marshal event to JSON

```go
jsonBytes, err := json.Marshal(event)
if err != nil {
    log.Fatal(err)
}
// Standard encoding/json works with Event struct tags
```

#### Unmarshal event from JSON

```go
var event events.Event
err := json.Unmarshal(jsonBytes, &event)
if err != nil {
    log.Fatal(err)
}

// Validate after unmarshaling
if err := events.Validate(event); err != nil {
    log.Printf("Received invalid event: %v", err)
}
```

---

### Filter Creation

#### Basic filter with standard fields

```go
since := int(time.Now().Add(-24 * time.Hour).Unix())
limit := 50

filter := filters.Filter{
    IDs:     []string{"abc123", "def456"},  // Prefix match
    Authors: []string{"cfa87f35"},          // Prefix match
    Kinds:   []int{1, 6, 7},
    Since:   &since,
    Limit:   &limit,
}
```

#### Filter with tag conditions

```go
filter := filters.Filter{
    Kinds: []int{1},
    Tags: filters.TagFilters{
        "e": {"5c83da77af1dec6d7289834998ad7aafbd9e2191396d75ec3cc27f5a77226f36"},
        "p": {"91cf9b32f3735070f46c0a86a820a47efa08a5be6c9f4f8cf68e5b5b75c92d60"},
    },
}
```

#### Filter with extensions (custom fields)

```go
// Extensions allow arbitrary JSON fields beyond the standard filter spec.
// For example, this is how to implement non-standard filters like 'search'.
filter := filters.Filter{
    Kinds: []int{1},
    Extensions: filters.FilterExtensions{
        "search": json.RawMessage(`"bitcoin"`),
    },
}

// Extensions are preserved during marshal/unmarshal but ignored by Matches().
// Storage/transport layers can inspect Extensions to implement custom behavior.
```

---

### Filter Matching

#### Match single event

```go
filter := filters.Filter{
    Authors: []string{"cfa87f35"},
    Kinds:   []int{1},
}

if filters.Matches(filter, event) {
    // Event satisfies all filter conditions
}
```

#### Filter event collection

```go
since := int(time.Now().Add(-1 * time.Hour).Unix())
filter := filters.Filter{
    Kinds: []int{1},
    Since: &since,
    Tags: filters.TagFilters{
        "p": {"abc123", "def456"},  // OR within tag values
    },
}

var matches []events.Event
for _, event := range events {
    if filters.Matches(filter, event) {
        matches = append(matches, event)
    }
}
```

---

### Filter JSON

#### Marshal filter to JSON

```go
filter := filters.Filter{
    IDs:   []string{"abc123"},
    Kinds: []int{1},
    Tags: filters.TagFilters{
        "e": {"event-id"},
    },
    Extensions: filters.FilterExtensions{
        "search": json.RawMessage(`"nostr"`),
    },
}

jsonBytes, err := filters.MarshalJSON(filter)
// Result: {"ids":["abc123"],"kinds":[1],"#e":["event-id"],"search":"nostr"}
```

#### Unmarshal filter from JSON

```go
jsonData := `{
    "authors": ["cfa87f35"],
    "kinds": [1],
    "#e": ["abc123"],
    "since": 1234567890,
    "search": "bitcoin"
}`

var filter filters.Filter
err := filters.UnmarshalJSON([]byte(jsonData), &filter)
if err != nil {
    log.Fatal(err)
}

// Standard fields populated: Authors, Kinds, Since
// Tag filters populated: Tags["e"] = ["abc123"]
// Unknown fields populated: Extensions["search"] = "bitcoin"
```

#### Extensions field behavior

The `Extensions` field captures any JSON properties not recognized as standard filter fields or tag filters. This design allows the core library to remain frozen while storage and transport layers implement custom filtering behavior.

**Standard fields**: `ids`, `authors`, `kinds`, `since`, `until`, `limit`

**Tag filters**: Any key starting with `#` (e.g., `#e`, `#p`, `#emoji`)

**Extensions**: Everything else

During marshaling, Extensions merge into the output JSON. During unmarshaling, unrecognized fields populate Extensions. The `Matches()` method ignores Extensions, and the library expects higher protocol layers to implement their usage.

Example implementing search filter:

```go
filter := filters.Filter{
    Kinds: []int{1},
    Extensions: filters.FilterExtensions{
        "search": json.RawMessage(`"bitcoin"`),
    },
}

// In a storage layer (not this library):
if searchRaw, ok := filter.Extensions["search"]; ok {
    var searchTerm string
    json.Unmarshal(searchRaw, &searchTerm)
    // Apply full-text search using searchTerm
}
```

## Testing

This library contains a comprehensive suite of unit tests. Run them with:

```bash
go test ./...
```
