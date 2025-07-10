package erx

type Coder interface {
	GetCode() string
	GetMsg() string
}

type Code string

const (
	Test Code = "test"
)

func (c Code) Code() string {
	return string(c)
}
