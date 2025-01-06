package erx

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/google/uuid"
)

type ErrorCallStack struct {
	CallStack   string
	OriginalErr error
}

func (e *ErrorCallStack) Error() string {
	return e.OriginalErr.Error()
}

func (e *ErrorCallStack) Unwrap() error {
	return e.OriginalErr
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}

// Wrap wraps the given error with a call stack. If the error is already an
// ErrorCallStack, it appends additional context (texts) to it. Otherwise,
// it converts the error to an ErrorCallStack and records the call stack.
func W(err error, texts ...string) error {
	if err == nil {
		return nil
	}

	var newErr error
	if len(texts) > 0 {
		errs := make([]error, len(texts))
		for i, text := range texts {
			errs[i] = errors.New(text)
		}
	}

	var errCStk *ErrorCallStack
	if errors.As(err, &errCStk) {
		return errors.Join(newErr, errCStk)
	}

	return errors.Join(newErr, &ErrorCallStack{
		CallStack:   getCallStack(3),
		OriginalErr: err,
	})
}

func New(text string) *ErrorCallStack {
	return &ErrorCallStack{
		CallStack:   getCallStack(2),
		OriginalErr: errors.New(text),
	}
}

func Errorf(format string, args ...interface{}) error {
	return New(fmt.Sprintf(format, args...))
}

func Log(err error) {
	var errCStk *ErrorCallStack
	if !errors.As(err, &errCStk) {
		return
	}

	id := uuid.New().String()

	cliMsg := fmt.Sprintf("%s\n%s", id, errCStk.OriginalErr.Error())
	FullMsg := fmt.Sprintf("%s\n%s\n%s", id, errCStk.OriginalErr.Error(), errCStk.CallStack)

	fmt.Println(cliMsg)
	fmt.Println(FullMsg)
}

// default use 2
func getCallStack(callerSkip ...int) (stkMsg string) {
	pc := make([]uintptr, 10)
	skip := 2
	if len(callerSkip) > 0 {
		skip = callerSkip[0]
	}
	n := runtime.Callers(skip, pc)

	frames := runtime.CallersFrames(pc[:n-2])

	stkMsg = "Call stack:\n"

	for {
		frame, more := frames.Next()
		stkMsg += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return stkMsg
}
