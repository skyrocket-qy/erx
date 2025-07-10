package erx

var MaxCallStackSize = 10

var ErrToCode = func(err error) Coder {
	return Unknown
}

type CoderImp string

const (
	Unknown CoderImp = "Unknown"
)

func (c CoderImp) Code() string {
	return string(c)
}

func (c CoderImp) Msg() string {
	return string(c)
}
