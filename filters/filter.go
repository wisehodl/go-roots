package filters

import (
	"encoding/json"
	"git.wisehodl.dev/jay/go-roots/events"
	"strings"
)

// TagFilters maps tag names to arrays of values for tag-based filtering
// Keys correspond to tag names without the "#" prefix.
type TagFilters map[string][]string

// FilterExtensions holds arbitrary additional filter fields as raw JSON.
// Allows custom filter extensions without modifying the core Filter type.
type FilterExtensions map[string]json.RawMessage

// Filter defines subscription criteria for events.
// All conditions within a filter applied with AND logic.
type Filter struct {
	IDs        []string
	Authors    []string
	Kinds      []int
	Since      *int
	Until      *int
	Limit      *int
	Tags       TagFilters
	Extensions FilterExtensions
}

// MarshalJSON converts the filter to JSON with standard fields, tag filters
// (prefixed with "#"), and extensions merged into a single object.
func MarshalJSON(f Filter) ([]byte, error) {
	outputMap := make(map[string]interface{})

	// Add standard fields
	if f.IDs != nil {
		outputMap["ids"] = f.IDs
	}
	if f.Authors != nil {
		outputMap["authors"] = f.Authors
	}
	if f.Kinds != nil {
		outputMap["kinds"] = f.Kinds
	}
	if f.Since != nil {
		outputMap["since"] = *f.Since
	}
	if f.Until != nil {
		outputMap["until"] = *f.Until
	}
	if f.Limit != nil {
		outputMap["limit"] = *f.Limit
	}

	// Add tags
	for key, values := range f.Tags {
		outputMap["#"+key] = values
	}

	// Merge extensions
	for key, raw := range f.Extensions {
		// Disallow standard keys in extensions
		if key == "ids" ||
			key == "authors" ||
			key == "kinds" ||
			key == "since" ||
			key == "until" ||
			key == "limit" {
			continue
		}

		// Disallow tag keys in extensions
		if strings.HasPrefix(key, "#") {
			continue
		}

		var extValue interface{}
		if err := json.Unmarshal(raw, &extValue); err != nil {
			return nil, err
		}
		outputMap[key] = extValue
	}

	return json.Marshal(outputMap)
}

// UnmarshalJSON parses JSON into the filter, separating standard fields,
// tag filters (keys starting with "#"), and extensions.
func UnmarshalJSON(data []byte, f *Filter) error {
	// Decode into raw map
	raw := make(FilterExtensions)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract standard fields
	if v, ok := raw["ids"]; ok {
		if err := json.Unmarshal(v, &f.IDs); err != nil {
			return err
		}
		delete(raw, "ids")
	}

	if v, ok := raw["authors"]; ok {
		if err := json.Unmarshal(v, &f.Authors); err != nil {
			return err
		}
		delete(raw, "authors")
	}

	if v, ok := raw["kinds"]; ok {
		if err := json.Unmarshal(v, &f.Kinds); err != nil {
			return err
		}
		delete(raw, "kinds")
	}

	if v, ok := raw["since"]; ok {
		if len(v) == 4 && string(v) == "null" {
			f.Since = nil
		} else {
			var val int
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}
			f.Since = &val
		}
		delete(raw, "since")
	}

	if v, ok := raw["until"]; ok {
		if len(v) == 4 && string(v) == "null" {
			f.Until = nil
		} else {
			var val int
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}
			f.Until = &val
		}
		delete(raw, "until")
	}

	if v, ok := raw["limit"]; ok {
		if len(v) == 4 && string(v) == "null" {
			f.Limit = nil
		} else {
			var val int
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}
			f.Limit = &val
		}
		delete(raw, "limit")
	}

	// Extract tag fields
	for key := range raw {
		if strings.HasPrefix(key, "#") {
			// Leave Tags as `nil` unless tag fields exist
			if f.Tags == nil {
				f.Tags = make(TagFilters)
			}
			tagKey := key[1:]
			var tagValues []string
			if err := json.Unmarshal(raw[key], &tagValues); err != nil {
				return err
			}
			f.Tags[tagKey] = tagValues
			delete(raw, key)
		}
	}

	// Place remaining fields in extensions
	if len(raw) > 0 {
		f.Extensions = raw
	}

	return nil
}

// Matches returns true if the event satisfies all filter conditions.
// Supports prefix matching for IDs and authors, and tag filtering.
// Does not account for custom extensions.
func Matches(f Filter, event events.Event) bool {
	// Check ID
	if len(f.IDs) > 0 {
		if !matchesPrefix(event.ID, f.IDs) {
			return false
		}
	}

	// Check Author
	if len(f.Authors) > 0 {
		if !matchesPrefix(event.PubKey, f.Authors) {
			return false
		}
	}

	// Check Kind
	if len(f.Kinds) > 0 {
		if !matchesKinds(event.Kind, f.Kinds) {
			return false
		}
	}

	// Check Timestamp
	if !matchesTimeRange(event.CreatedAt, f.Since, f.Until) {
		return false
	}

	// Check Tags
	if len(f.Tags) > 0 {
		if !matchesTags(event.Tags, &f.Tags) {
			return false
		}
	}

	return true
}

func matchesPrefix(candidate string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(candidate, prefix) {
			return true
		}
	}
	return false
}

func matchesKinds(candidate int, kinds []int) bool {
	for _, kind := range kinds {
		if candidate == kind {
			return true
		}
	}
	return false
}

func matchesTimeRange(timestamp int, since *int, until *int) bool {
	if since != nil && timestamp < *since {
		return false
	}
	if until != nil && timestamp > *until {
		return false
	}
	return true
}

func matchesTags(eventTags []events.Tag, tagFilters *TagFilters) bool {
	// Build index of tags and values
	eventIndex := make(map[string][]string, len(eventTags))
	for _, tag := range eventTags {
		if len(tag) < 2 {
			continue
		}
		eventIndex[tag[0]] = append(eventIndex[tag[0]], tag[1])
	}

	// Check filters against the index
	for tagName, filterValues := range *tagFilters {
		// Skip empty tag filters (empty tag filters match all events)
		if len(filterValues) == 0 {
			continue
		}

		eventValues, exists := eventIndex[tagName]
		if !exists {
			return false
		}

		found := false
		for _, filterVal := range filterValues {
			for _, eventVal := range eventValues {
				if eventVal == filterVal {
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			return false
		}
	}

	// If no filter explicitly fails, then the event is matched
	return true
}
