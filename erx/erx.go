package erx

import (
	"errors"
	"fmt"
	"runtime"
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
	var a, b error
	e1, ok := err.(*ErrorCallStack)
	if ok {
		a = e1.OriginalErr
	} else {
		a = err
	}
	e2, ok := target.(*ErrorCallStack)
	if ok {
		b = e2.OriginalErr
	} else {
		b = target
	}

	return errors.Is(a, b)
}

func As(err error, target any) bool {
	ce, ok := err.(*ErrorCallStack)
	if ok {
		return errors.As(ce.OriginalErr, target)
	}

	return errors.As(err, target)
}

func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var err error
	for _, e := range errs {
		if e == nil {
			continue
		}
		e1, ok := e.(*ErrorCallStack)
		if ok {
			err = errors.Join(err, e1.OriginalErr)
		} else {
			err = errors.Join(err, e)
		}
	}
	return err
}

// Check err is original error, if true, convert to ErrorCallStack which record
// call stack
func Wrap(msg string, err error) error {
	e1, ok := err.(*ErrorCallStack)
	if ok {
		return Join(errors.New(msg), e1)
	}
	return Join(errors.New(msg), &ErrorCallStack{
		CallStack:   getCallStack(),
		OriginalErr: err,
	})
}

func New(text string) *ErrorCallStack {
	return &ErrorCallStack{
		CallStack:   getCallStack(),
		OriginalErr: errors.New(text),
	}
}

func Errorf(format string, args ...interface{}) error {
	return New(fmt.Sprintf(format, args...))
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

	stkMsg = "Call stack:"

	for {
		frame, more := frames.Next()
		stkMsg += fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
	return stkMsg
}
