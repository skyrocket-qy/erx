package erx

import (
	"errors"
	"testing"
)

func TestDefaultErrToCodeCoverage(t *testing.T) {
	stdErr := errors.New("any error")
	code := ErrToCode(stdErr)
	if code != ErrUnknown {
		t.Errorf("expected default code to be ErrUnknown, got %v", code)
	}
}
