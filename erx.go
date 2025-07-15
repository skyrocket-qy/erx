package erx

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

// W wraps the given error with a call stack and optional additional context.
//
// If the error is already a contextError, it prepends the first optional text message
// to the existing error using errors.Join. The call stack is not updated in this case.
//
// If the error is not a contextError, it creates a new contextError with:
//   - The original error
//   - A generated call stack (skipping 3 frames)
//   - An associated error code via ErrToCode
//
// Note: Only the first element of texts (if any) is used. texts[1:] are ignored.
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

func WCode(code Coder, err error, texts ...string) error {
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

// New creates a new contextError with the given Coder and optional message.
//
// It captures a call stack (skipping 2 frames) and assigns the provided error code.
// The final error message is constructed from the Coder's message,
// optionally appending the first element of texts if provided.
//
// Note: Only texts[0] is used if texts is non-empty. Additional elements are ignored.
func New(coder Coder, texts ...string) *contextError {
	ctxErr := &contextError{
		callerInfos: getCallStack(2),
		code:        coder,
	}

	errMsg := coder.Code()
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

func GetClientMsg(err error) (code Coder, ok bool) {
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

func (im InternalMsg) MarshalZerologObject(e *zerolog.Event) {
	if im.Cause != nil {
		e.Str("cause", im.Cause.Error())
	}
	if im.Err != nil {
		e.Str("error", im.Err.Error())
	}
	e.Str("code", im.Code)

	// Convert callerInfos to pretty strings
	filtered := filterCallerInfos(im.CallerInfos)
	trace := make([]string, 0, len(filtered))
	for _, ci := range filtered {
		trace = append(trace, fmt.Sprintf("%s %d %s",
			trimToProject(ci.File),
			ci.Line,
			extractFuncName(ci.Function),
		))
	}

	// formatted := make([]string, 0, len(im.CallerInfos))
	// for _, ci := range im.CallerInfos {
	// 	shortPath := trimToProject(ci.File)
	// 	shortFunc := extractFuncName(ci.Function)
	// 	formatted = append(formatted, fmt.Sprintf("%s %d %s", shortPath, ci.Line, shortFunc))
	// }
	e.Strs("callerTrace", trace)
}

func trimToProject(path string) string {
	const projectRoot = "/Users/qy/skyro/author/"
	if rel, ok := strings.CutPrefix(path, projectRoot); ok {
		return rel
	}
	return path
}

func extractFuncName(fullFunc string) string {
	// e.g., input: srv/internal/logic/inter.(*Logic).Login
	// output: (*Logic).Login
	if idx := strings.LastIndex(fullFunc, "/"); idx >= 0 {
		return fullFunc[idx+1:]
	}
	return fullFunc
}

func filterCallerInfos(infos []CallerInfo) []CallerInfo {
	const projectPrefix = "/Users/qy/skyro/author/"
	var filtered []CallerInfo
	for _, ci := range infos {
		if strings.HasPrefix(ci.File, projectPrefix) {
			filtered = append(filtered, ci)
		}
	}
	return filtered
}
