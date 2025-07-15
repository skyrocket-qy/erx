package erx

type Coder interface {
	Code() string
}

var _ Coder = CoderImp("")

type CoderImp string

const (
	ErrUnknown CoderImp = "500.0000"
)

func (c CoderImp) Code() string {
	return string(c)
}
