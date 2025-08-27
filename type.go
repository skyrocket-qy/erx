package erx

type CallerInfo struct {
	Function string
	File     string
	Line     int
	Msg      string
}

type CtxErr struct {
	Code        Code
	CallerInfos []CallerInfo
	cause       error // original error
}

func (e *CtxErr) Error() string {
	return e.Code.Str()
}

func (e *CtxErr) Unwrap() error {
	return e.cause
}

func (e *CtxErr) SetCode(c Code) *CtxErr {
	if e == nil {
		return nil
	}

	e.Code = c
	return e
}
