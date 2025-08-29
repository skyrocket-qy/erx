package erx_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/skyrocket-qy/erx"
)

func TestMain(m *testing.M) {
	originalErrToCode := erx.ErrToCode
	erx.ErrToCode = func(err error) erx.Code {
		if strings.Contains(err.Error(), "a custom error for mapping") {
			return ErrCustom
		}
		return erx.ErrUnknown
	}

	code := m.Run()

	erx.ErrToCode = originalErrToCode
	os.Exit(code)
}

const (
	ErrTest erx.CodeImp = "500.0001"
)

func TestNew(t *testing.T) {
	err := erx.New(ErrTest, "test message")
	if err == nil {
		t.Fatal("erx.New should not return nil")
	}

	if err.Code != ErrTest {
		t.Errorf("expected code %v, got %v", ErrTest, err.Code)
	}

	if len(err.CallerInfos) == 0 {
		t.Fatal("expected caller info, got none")
	}

	if !strings.Contains(err.CallerInfos[0].Msg, "test message") {
		t.Errorf("expected message 'test message', got '%s'", err.CallerInfos[0].Msg)
	}
}

func TestW_StandardError(t *testing.T) {
	stdErr := errors.New("a standard error")
	err := erx.W(stdErr, "wrapped")

	if err == nil {
		t.Fatal("erx.W should not return nil")
	}

	if !errors.Is(err, stdErr) {
		t.Fatal("expected error to be wrapp`ing stdErr")
	}

	if err.Code != erx.ErrUnknown {
		t.Errorf("expected code %v, got %v", erx.ErrUnknown, err.Code)
	}

	if len(err.CallerInfos) == 0 {
		t.Fatal("expected caller info, got none")
	}

	if !strings.Contains(err.CallerInfos[0].Msg, "wrapped") {
		t.Errorf("expected message 'wrapped', got '%s'", err.CallerInfos[0].Msg)
	}
}

func TestWf(t *testing.T) {
	stdErr := errors.New("a standard error")
	err := erx.Wf(stdErr, "wrapped %s", "message")

	if err == nil {
		t.Fatal("erx.Wf should not return nil")
	}

	if !errors.Is(err, stdErr) {
		t.Fatal("expected error to be wrapping stdErr")
	}

	if err.Code != erx.ErrUnknown {
		t.Errorf("expected code %v, got %v", erx.ErrUnknown, err.Code)
	}

	if len(err.CallerInfos) == 0 {
		t.Fatal("expected caller info, got none")
	}

	if !strings.Contains(err.CallerInfos[0].Msg, "wrapped message") {
		t.Errorf("expected message 'wrapped message', got '%s'", err.CallerInfos[0].Msg)
	}
}

func TestW_NilError(t *testing.T) {
	err := erx.W(nil)
	if err != nil {
		t.Fatal("erx.W(nil) should return nil")
	}
}

func TestErrorsAs(t *testing.T) {
	err := erx.New(ErrTest, "test message")
	var erxErr *erx.CtxErr
	if !errors.As(err, &erxErr) {
		t.Fatal("errors.As should be able to extract CtxErr")
	}
	if erxErr.Code != ErrTest {
		t.Errorf("extracted error has wrong code: got %v, want %v", erxErr.Code, ErrTest)
	}
}

func TestW_ErxError(t *testing.T) {
	erxErr := erx.New(ErrTest, "initial error")
	err := erx.W(erxErr, "wrapped context")

	if err != erxErr {
		t.Fatal("erx.W should return the same error instance")
	}

	if err.Code != ErrTest {
		t.Errorf("expected code %v, got %v", ErrTest, err.Code)
	}
}

const (
	ErrCustom erx.CodeImp = "404.0001"
)

func TestErrToCode(t *testing.T) {
	stdErr := errors.New("a custom error for mapping")
	err := erx.W(stdErr, "wrapped")

	if err.Code != ErrCustom {
		t.Errorf("expected code %v, got %v", ErrCustom, err.Code)
	}
}

func TestCallStack(t *testing.T) {
	err := erx.New(ErrTest, "test message")

	if len(err.CallerInfos) < 1 {
		t.Fatal("expected at least one caller info")
	}

	caller := err.CallerInfos[0]
	if !strings.Contains(caller.Function, "erx_test.TestCallStack") {
		t.Errorf("expected function name 'erx_test.TestCallStack', got '%s'", caller.Function)
	}
	if !strings.HasSuffix(caller.File, "erx_test.go") {
		t.Errorf("expected file name 'erx_test.go', got '%s'", caller.File)
	}
	if caller.Line == 0 {
		t.Error("expected line number, got 0")
	}
}

func TestNewf(t *testing.T) {
	err := erx.Newf(ErrTest, "test message %d", 1)
	if err == nil {
		t.Fatal("erx.Newf should not return nil")
	}

	if err.Code != ErrTest {
		t.Errorf("expected code %v, got %v", ErrTest, err.Code)
	}

	if len(err.CallerInfos) == 0 {
		t.Fatal("expected caller info, got none")
	}

	if !strings.Contains(err.CallerInfos[0].Msg, "test message 1") {
		t.Errorf("expected message 'test message 1', got '%s'", err.CallerInfos[0].Msg)
	}
}

func TestFullMsg(t *testing.T) {
	// Case 1: nil error
	if erx.FullMsg(nil) != "" {
		t.Error("FullMsg(nil) should be empty string")
	}

	// Case 2: standard error
	stdErr := errors.New("a standard error")
	if erx.FullMsg(stdErr) != "a standard error" {
		t.Errorf("FullMsg on standard error failed, got '%s'", erx.FullMsg(stdErr))
	}

	// Case 3: CtxErr with no cause
	err := erx.New(ErrTest, "test message")
	fullMsg := erx.FullMsg(err)
	if !strings.Contains(fullMsg, ErrTest.Str()) {
		t.Errorf("FullMsg should contain error code, got '%s'", fullMsg)
	}
	if !strings.Contains(fullMsg, "erx_test.TestFullMsg") {
		t.Errorf("FullMsg should contain caller info, got '%s'", fullMsg)
	}
	if !strings.Contains(fullMsg, "test message") {
		t.Errorf("FullMsg should contain message, got '%s'", fullMsg)
	}

	// Case 4: CtxErr with a cause
	errWithCause := erx.W(stdErr, "wrapped")
	fullMsgWithCause := erx.FullMsg(errWithCause)
	if !strings.Contains(fullMsgWithCause, erx.ErrUnknown.Str()) {
		t.Errorf("FullMsg with cause should contain error code, got '%s'", fullMsgWithCause)
	}
	if !strings.Contains(fullMsgWithCause, "wrapped") {
		t.Errorf("FullMsg with cause should contain message, got '%s'", fullMsgWithCause)
	}
	if !strings.Contains(fullMsgWithCause, "Caused by: a standard error") {
		t.Errorf("FullMsg with cause should contain cause, got '%s'", fullMsgWithCause)
	}
}

func TestCtxErrMethods(t *testing.T) {
	// Test Unwrap
	stdErr := errors.New("a standard error")
	errWithCause := erx.W(stdErr, "wrapped")
	if !errors.Is(errWithCause, stdErr) {
		t.Fatal("errors.Is should work with unwrapped error")
	}
	if errors.Unwrap(errWithCause) != stdErr {
		t.Fatal("errors.Unwrap should return the cause")
	}

	errWithoutCause := erx.New(ErrTest, "test")
	if errors.Unwrap(errWithoutCause) != nil {
		t.Fatal("errors.Unwrap on error from New should be nil")
	}

	// Test SetCode
	err := erx.New(ErrTest, "test")
	const newCode erx.CodeImp = "123.456"
	err.SetCode(newCode)
	if err.Code != newCode {
		t.Errorf("SetCode failed, expected code %v, got %v", newCode, err.Code)
	}

	// Test SetCode on nil receiver
	var nilErr *erx.CtxErr
	if nilErr.SetCode(newCode) != nil {
		t.Error("SetCode on nil receiver should return nil")
	}

	// Test Error()
	if err.Error() != newCode.Str() {
		t.Errorf("Error() should return code string, got '%s'", err.Error())
	}
}
