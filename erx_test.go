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
