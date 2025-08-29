package erx

import (
	"errors"
	"fmt"
	"runtime"
)

// W wraps the given error with a call stack and optional additional context.
func W(err error, msgs ...string) *CtxErr {
	return w(err, msgs...)
}

func Wf(err error, format string, args ...any) *CtxErr {
	return w(err, fmt.Sprintf(format, args...))
}

func w(err error, msgs ...string) *CtxErr {
	if err == nil {
		return nil
	}

	var ctxErr *CtxErr
	if !errors.As(err, &ctxErr) {
		ctxErr = &CtxErr{
			cause:       err,
			CallerInfos: getCallStack(4),
			Code:        ErrToCode(err),
		}
	}

	if len(msgs) > 0 && len(ctxErr.CallerInfos) > 0 {
		ctxErr.CallerInfos[0].Msg += " " + msgs[0]
	}

	return ctxErr
}

func New(code Code, msgs ...string) *CtxErr {
	ctxErr := &CtxErr{
		CallerInfos: getCallStack(3),
		Code:        code,
	}

	if len(msgs) > 0 {
		ctxErr.CallerInfos[0].Msg = msgs[0]
	}

	return ctxErr
}

func Newf(code Code, format string, args ...any) *CtxErr {
	ctxErr := &CtxErr{
		CallerInfos: getCallStack(3),
		Code:        code,
	}

	if format != "" {
		ctxErr.CallerInfos[0].Msg = fmt.Sprintf(format, args...)
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

func FullMsg(err error) string {
	if err == nil {
		return ""
	}

	var ctxErr *CtxErr
	if !errors.As(err, &ctxErr) {
		return err.Error()
	}

	// Start with error code string
	msg := ctxErr.Code.Str()

	// Add caller infos if available
	if len(ctxErr.CallerInfos) > 0 {
		for _, ci := range ctxErr.CallerInfos {
			msg += fmt.Sprintf("\n  at %s (%s:%d)", ci.Function, ci.File, ci.Line)
			if ci.Msg != "" {
				msg += fmt.Sprintf(": %s", ci.Msg)
			}
		}
	}

	// Recurse into cause if any
	if ctxErr.cause != nil {
		msg += "\nCaused by: " + FullMsg(ctxErr.cause)
	}

	return msg
}
