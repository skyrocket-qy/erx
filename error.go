package erx

type callerInfo struct {
	Function string
	File     string
	Line     int
}

type contextError struct {
	err         error
	code        Coder
	callerInfos []callerInfo
}

func (e *contextError) Error() string {
	return e.code.GetMsg()
}

func (e *contextError) getCallerInfos() []callerInfo {
	return e.callerInfos
}
