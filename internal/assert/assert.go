package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, actual, expected T) {

	// to annotate the go test runner that this function
	// is a helper function
	t.Helper()
	if actual != expected {
		t.Errorf("got %v; want %v", actual, actual)
	}
}
