package roots

import (
	"encoding/json"
	// "fmt"
	"strings"
)

type Filter struct {
	IDs        []string
	Authors    []string
	Kinds      []int
	Since      *int
	Until      *int
	Limit      *int
	Tags       map[string][]string
	Extensions map[string]json.RawMessage
}

func (f *Filter) MarshalJSON() ([]byte, error) {
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

func (f *Filter) UnmarshalJSON(data []byte) error {
	// Decode into raw map
	raw := make(map[string]json.RawMessage)
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
				f.Tags = make(map[string][]string)
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

func (f *Filter) Matches(event *Event) bool {
	// Check ID
	if len(f.IDs) > 0 {
		if !matchesPrefix(event.ID, &f.IDs) {
			return false
		}
	}

	// Check Author
	if len(f.Authors) > 0 {
		if !matchesPrefix(event.PubKey, &f.Authors) {
			return false
		}
	}

	// Check Kind
	if len(f.Kinds) > 0 {
		if !matchesKinds(event.Kind, &f.Kinds) {
			return false
		}
	}

	// Check Timestamp
	if !matchesTimeRange(event.CreatedAt, f.Since, f.Until) {
		return false
	}

	// Check Tags
	if len(f.Tags) > 0 {
		if !matchesTags(&event.Tags, &f.Tags) {
			return false
		}
	}

	return true
}

func matchesPrefix(candidate string, prefixes *[]string) bool {
	for _, prefix := range *prefixes {
		if strings.HasPrefix(candidate, prefix) {
			return true
		}
	}
	return false
}

func matchesKinds(candidate int, kinds *[]int) bool {
	for _, kind := range *kinds {
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

func matchesTags(eventTags *[][]string, filterTags *map[string][]string) bool {
	for tagName, filterValues := range *filterTags {
		if len(filterValues) == 0 {
			return true
		}

		found := false
		for _, eventTag := range *eventTags {
			if len(eventTag) < 2 {
				continue
			}
			if eventTag[0] != tagName {
				continue
			}

			for _, filterValue := range filterValues {
				if eventTag[1] == filterValue {
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

	return true
}
