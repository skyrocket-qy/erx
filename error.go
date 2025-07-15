package erx

type CallerInfo struct {
	Function string
	File     string
	Line     int
}

type CxtErr struct {
	ThirdpartyErr *error
	Code          Code
	CallerInfos   []CallerInfo
}

func (e *CxtErr) Error() string {
	return e.Code.Str()
}
