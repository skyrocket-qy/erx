package erx_test

import (
	"errors"
	"testing"

	"github.com/skyrocketOoO/erx/erx"
)

func TestWrapPreservesOriginalError(t *testing.T) {
	// Define an original error
	originalErr := errors.New("original error")

	// Wrap the error without additional context
	wrappedErr := erx.W(originalErr)

	// Assert that the original error is preserved
	if !errors.Is(wrappedErr, originalErr) {
		t.Errorf("wrapped error does not match original error: got %v, want %v", wrappedErr, originalErr)
	}

	// Extract ErrorCtx and check the OriginalErr
	var errCtx *erx.ErrorCtx
	if !errors.As(wrappedErr, &errCtx) {
		t.Errorf("wrapped error is not of type *ErrorCtx: got %T", wrappedErr)
	}
}

func TestWrapErrorWithAdditionalContext(t *testing.T) {
	// Define an original error
	originalErr := errors.New("original error")

	// Wrap the error with additional context
	texts := []string{"context 1", "context 2"}
	wrappedErr := erx.W(originalErr, texts...)

	// Assert that the original error is preserved
	if !errors.Is(wrappedErr, originalErr) {
		t.Errorf("wrapped error does not match original error: got %v, want %v", wrappedErr, originalErr)
	}

	// Extract ErrorCtx and check the OriginalErr
	var errCtx *erx.ErrorCtx
	if !errors.As(wrappedErr, &errCtx) {
		t.Errorf("wrapped error is not of type *ErrorCtx: got %T", wrappedErr)
	}

	if !errors.Is(errCtx.OriginalErr, originalErr) {
		t.Errorf("OriginalErr in ErrorCtx does not match original error: got %v, want %v", errCtx.OriginalErr, originalErr)
	}
}

func TestNewErrorCtx(t *testing.T) {
	// Create a new ErrorCtx using New
	text := "new error"
	errCtx := erx.New(text)

	// Check that the OriginalErr is set correctly
	if errCtx.Error() != text {
		t.Errorf("OriginalErr in ErrorCtx does not match input text: got %v, want %v", errCtx.OriginalErr.Error(), text)
	}

	// Check that the CallStack and ID are populated
	if errCtx.Ctx["CallStack"] == "" {
		t.Errorf("CallStack in ErrorCtx is empty")
	}
	if errCtx.Ctx["ID"] == "" {
		t.Errorf("ID in ErrorCtx is empty")
	}
}

func TestUnwrapPreservesOriginalError(t *testing.T) {
	// Define an original error
	originalErr := errors.New("original error")

	// Wrap the error
	wrappedErr := erx.W(originalErr)

	// Extract ErrorCtx and unwrap the error
	var errCtx *erx.ErrorCtx
	if !errors.As(wrappedErr, &errCtx) {
		t.Errorf("wrapped error is not of type *ErrorCtx: got %T", wrappedErr)
	}

	if unwrappedErr := errCtx.Unwrap(); !errors.Is(unwrappedErr, originalErr) {
		t.Errorf("Unwrap did not return the original error: got %v, want %v", unwrappedErr, originalErr)
	}
}
