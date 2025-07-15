package erx

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
	return e.code.Code()
}

func (e *contextError) getCallerInfos() []CallerInfo {
	return e.callerInfos
}
