package roots

import (
	"strings"
	"testing"
)

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
