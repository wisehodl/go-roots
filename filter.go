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

func (filter Filter) MarshalJSON() ([]byte, error) {
	outputMap := make(map[string]interface{})

	// Add standard fields
	if filter.IDs != nil {
		outputMap["ids"] = filter.IDs
	}
	if filter.Authors != nil {
		outputMap["authors"] = filter.Authors
	}
	if filter.Kinds != nil {
		outputMap["kinds"] = filter.Kinds
	}
	if filter.Since != nil {
		outputMap["since"] = *filter.Since
	}
	if filter.Until != nil {
		outputMap["until"] = *filter.Until
	}
	if filter.Limit != nil {
		outputMap["limit"] = *filter.Limit
	}

	// Add tags
	for key, values := range filter.Tags {
		outputMap["#"+key] = values
	}

	// Merge extensions
	for key, raw := range filter.Extensions {
		// Disallow standard keys in extensions
		standardKeys := []string{"ids", "authors", "kinds", "since", "until", "limit"}
		found := false
		for _, stdKey := range standardKeys {
			// Skip standard key
			if key == stdKey {
				found = true
				break
			}
		}
		if found {
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

func (filter *Filter) UnmarshalJSON(data []byte) error {
	// Decode into raw map
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Extract standard fields
	if v, ok := raw["ids"]; ok {
		if err := json.Unmarshal(v, &filter.IDs); err != nil {
			return err
		}
		delete(raw, "ids")
	}

	if v, ok := raw["authors"]; ok {
		if err := json.Unmarshal(v, &filter.Authors); err != nil {
			return err
		}
		delete(raw, "authors")
	}

	if v, ok := raw["kinds"]; ok {
		if err := json.Unmarshal(v, &filter.Kinds); err != nil {
			return err
		}
		delete(raw, "kinds")
	}

	if v, ok := raw["since"]; ok {
		if string(raw["since"]) == "null" {
			filter.Since = nil
		} else {
			var val int
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}
			filter.Since = &val
		}
		delete(raw, "since")
	}

	if v, ok := raw["until"]; ok {
		if string(raw["until"]) == "null" {
			filter.Until = nil
		} else {
			var val int
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}
			filter.Until = &val
		}
		delete(raw, "until")
	}

	if v, ok := raw["limit"]; ok {
		if string(raw["limit"]) == "null" {
			filter.Limit = nil
		} else {
			var val int
			if err := json.Unmarshal(v, &val); err != nil {
				return err
			}
			filter.Limit = &val
		}
		delete(raw, "limit")
	}

	// Extract tag fields
	for key := range raw {
		if strings.HasPrefix(key, "#") {
			// Leave Tags as `nil` unless tag fields exist
			if filter.Tags == nil {
				filter.Tags = make(map[string][]string)
			}
			tagKey := strings.TrimPrefix(key, "#")
			var tagValues []string
			if err := json.Unmarshal(raw[key], &tagValues); err != nil {
				return err
			}
			filter.Tags[tagKey] = tagValues
			delete(raw, key)
		}
	}

	// Place remaining fields in extensions
	if len(raw) > 0 {
		filter.Extensions = raw
	}

	return nil
}

func (filter Filter) Matches(event Event) bool {
	// Check ID
	if len(filter.IDs) > 0 {
		if !matchesPrefix(event.ID, filter.IDs) {
			return false
		}
	}

	// Check Author
	if len(filter.Authors) > 0 {
		if !matchesPrefix(event.PubKey, filter.Authors) {
			return false
		}
	}

	// Check Kind
	if len(filter.Kinds) > 0 {
		if !matchesKinds(event.Kind, filter.Kinds) {
			return false
		}
	}

	// Check Timestamp
	if !matchesTimeRange(event.CreatedAt, filter.Since, filter.Until) {
		return false
	}

	// Check Tags
	if len(filter.Tags) > 0 {
		if !matchesTags(event.Tags, filter.Tags) {
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

func matchesTags(eventTags [][]string, filterTags map[string][]string) bool {
	for tagName, filterValues := range filterTags {
		if len(filterValues) == 0 {
			return true
		}

		found := false
		for _, eventTag := range eventTags {
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
