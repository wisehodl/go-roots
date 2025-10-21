package roots

import (
	"strings"
	"testing"
)

func intPtr(i int) *int {
	return &i
}

func expectOk(t *testing.T, err error) {
	if err != nil {
		t.Errorf("got error: %s", err.Error())
	}
}

func expectError(t *testing.T, err error) {
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func expectErrorSubstring(t *testing.T, err error, expected string) {
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("error = %q, want substring %q", err.Error(), expected)
	}
}

func expectEqualStrings(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func expectEqualIntPointers(t *testing.T, got, want *int) {
	if *got != *want {
		t.Errorf("got %d, want %d", *got, *want)
	}
}

func expectEqualStringSlices(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Errorf("length mismatch: got %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("index %d: got %s, want %s", i, got[i], want[i])
		}
	}
}

func expectEqualIntSlices(t *testing.T, got, want []int) {
	if len(got) != len(want) {
		t.Errorf("length mismatch: got %d, want %d", len(got), len(want))
		return
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("index %d: got %d, want %d", i, got[i], want[i])
		}
	}

}

func expectEqualEvents(t *testing.T, got, want Event) {
	if got.ID != want.ID ||
		got.PubKey != want.PubKey ||
		got.CreatedAt != want.CreatedAt ||
		got.Kind != want.Kind ||
		got.Content != want.Content ||
		got.Sig != want.Sig ||
		!equalTags(got.Tags, want.Tags) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func equalTags(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
