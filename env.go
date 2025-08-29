package erx

var MaxCallStackSize = 10

func defaultErrToCode(err error) Code {
	return ErrUnknown
}

var ErrToCode = defaultErrToCode
