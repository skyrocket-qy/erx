package erx

// type ErrorCtx struct {
// 	OriginalErr error
// 	Ctx         map[string]string
// }

// func (e *ErrorCtx) Error() string {
// 	return e.OriginalErr.Error()
// }

// func (e *ErrorCtx) Unwrap() error {
// 	return e.OriginalErr
// }

type CallerInfo struct {
	Function string
	File     string
	Line     int
}

type contextError struct {
	err         error
	code        Coder
	callerInfos []CallerInfo
}

func (e *contextError) Error() string {
	return e.err.Error()
}

func (e *contextError) getCallStack() []CallerInfo {
	return e.callerInfos
}
