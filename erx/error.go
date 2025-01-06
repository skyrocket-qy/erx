package erx

type ErrorCtx struct {
	OriginalErr error
	Ctx         map[string]string
}

func (e *ErrorCtx) Error() string {
	return e.OriginalErr.Error()
}

func (e *ErrorCtx) Unwrap() error {
	return e.OriginalErr
}
