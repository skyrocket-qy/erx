package erx

import (
	"errors"
	"runtime"
)

var (
	// TODO: Consider adding a default ErrToCode to avoid nil dereference if not set.
	// Provide helper functions for users to easily set ErrToCode and create errors.
	ErrToCode        func(err error) Coder
	MaxCallStackSize = 10
)

// Wraps the given error with a call stack. If the error is already an
// contextError, it appends additional context (texts) to it. Otherwise,
// it converts the error to an contextError and records the call stack.
// Not handle texts[1:]
func W(err error, texts ...string) error {
	if err == nil {
		return nil
	}

	var ctxErr *contextError
	if errors.As(err, &ctxErr) {
		if len(texts) > 0 {
			ctxErr.err = errors.Join(errors.New(texts[0]), ctxErr.err)
		}
	} else {
		ctxErr = &contextError{
			callerInfos: getCallStack(3),
			code:        ErrToCode(err),
			err:         err,
		}
	}

	return ctxErr
}

func New(coder Coder, texts ...string) *contextError {
	ctxErr := &contextError{
		callerInfos: getCallStack(2),
		code:        coder,
	}

	errMsg := coder.GetMsg()
	if len(texts) > 0 {
		errMsg = errMsg + ", " + texts[0]
	}

	ctxErr.err = errors.New(errMsg)
	return ctxErr
}

func getCallStack(callerSkip ...int) (callerInfos []CallerInfo) {
	pc := make([]uintptr, MaxCallStackSize)
	skip := 2
	if len(callerSkip) > 0 {
		skip = callerSkip[0]
	}
	n := runtime.Callers(skip, pc)

	frames := runtime.CallersFrames(pc[:n])

	for {
		frame, more := frames.Next()
		callerInfos = append(callerInfos, CallerInfo{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}

	return callerInfos
}

// Get the root cause of the error
// only support join erros or erx error
func Cause(err error) error {
	ctxErr, ok := err.(*contextError)
	if !ok {
		return err
	}

	err = ctxErr.err
	// two case: 1. join error 2. single error
	if jErr, ok := err.(interface{ Unwrap() []error }); ok {
		joinErrs := jErr.Unwrap()
		if len(joinErrs) > 0 {
			return joinErrs[len(joinErrs)-1]
		}
		return err
	}

	return err
}

func GetClientMsg(err error) (code string, ok bool) {
	var ctxErr *contextError
	if !errors.As(err, &ctxErr) {
		return "", false
	}

	return ctxErr.code.GetCode(), true
}

type InternalMsg struct {
	Cause       error
	Code        string
	CallerInfos []CallerInfo
	Err         error
}

func GetInternalMsg(err error) (msg InternalMsg, ok bool) {
	var ctxErr *contextError
	if !errors.As(err, &ctxErr) {
		return InternalMsg{}, false
	}

	return InternalMsg{
		Cause:       Cause(err),
		Code:        ctxErr.code.GetCode(),
		CallerInfos: ctxErr.getCallStack(),
		Err:         ctxErr.err,
	}, true
}
