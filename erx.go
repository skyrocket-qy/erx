package erx

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// W wraps the given error with a call stack and optional additional context.
// the ctx is tricky format, len = 0 : no context, len = 1: raw string, len > 1:
func W(err error, ctx ...string) error {
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

func WCode(code Code, err error, texts ...string) error {
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
			code:        code,
			err:         err,
		}
	}

	return ctxErr
}

func New(code Code, texts ...string) *CxtErr {
	ctxErr := &CxtErr{
		CallerInfos: getCallStack(2),
		Code:        code,
	}

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

// Cause extracts the root cause from a contextError.
//
// If the provided error is not a contextError, it returns the error as-is.
// If it is a contextError and wraps a joined error (via errors.Join),
// it returns the last error in the joined chain (assumed to be the original cause).
// Otherwise, it returns the wrapped error directly.
//
// Note: This function only unwraps one level. It does not recurse through nested joins.
func Cause(err error) error {
	ctxErr := &contextError{}

	ok := errors.As(err, &ctxErr)
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

func GetClientMsg(err error) (code Code, ok bool) {
	var ctxErr *contextError
	if !errors.As(err, &ctxErr) {
		return nil, false
	}

	return ctxErr.code, true
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
		Code:        ctxErr.code.Code(),
		CallerInfos: ctxErr.getCallerInfos(),
		Err:         ctxErr.err,
	}, true
}

func W(err error, ctx ...any) error {
	im, ok := err.(*InternalMsg)
	if !ok {
		im = New(err) // capture stack
	}

	// Add context to the correct frame
	file, line := currentFileLine()
	frame := findFrame(im.StackTrace, file, line)
	if frame == -1 {
		frame = 0 // fallback
	}

	// Process context
	switch len(ctx) {
	case 0:
		// no context
	case 1:
		// raw string
		if s, ok := ctx[0].(string); ok {
			im.PerFrame[frame] = append(im.PerFrame[frame], s)
		}
	default:
		// maybe key-value or format string
		if format, ok := ctx[0].(string); ok && strings.Contains(format, "%") {
			// treat as fmt.Sprintf
			im.PerFrame[frame] = append(im.PerFrame[frame], fmt.Sprintf(format, ctx[1:]...))
		} else if len(ctx)%2 == 0 {
			// treat as k/v
			for i := 0; i < len(ctx); i += 2 {
				k, ok1 := ctx[i].(string)
				v := fmt.Sprint(ctx[i+1])
				if !ok1 {
					continue
				}
				kv := fmt.Sprintf("%s:%s", k, v)
				im.PerFrame[frame] = append(im.PerFrame[frame], kv)
				// optionally: store in Fields map for logging
			}
		} else {
			// fallback: stringify everything
			joined := make([]string, len(ctx))
			for i, v := range ctx {
				joined[i] = fmt.Sprint(v)
			}
			im.PerFrame[frame] = append(im.PerFrame[frame], strings.Join(joined, " "))
		}
	}

	return im
}
