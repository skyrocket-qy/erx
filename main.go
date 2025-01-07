package main

import (
	"errors"
	"fmt"

	"github.com/skyrocketOoO/erx/erx"
)

var (
	ErrUnknown = errors.New("unknown error")
	ErrDb      = errors.New("db error")
)

var OgErrDb = errors.New("db error")

func main() {
	// es := errors.Join(ErrDb, ErrUnknown)
	// Use errors.As to extract the joined errors
	// fmt.Println(erx.Cause(ErrDb))

	// fmt.Println(errors.Is(erx.Errorf("%w", ErrDb), ErrDb))
	// fmt.Println(errors.Is(erx.W(ErrUnknown, "db error"), ErrUnknown))
	// fmt.Println(errors.Is(errors.Join(ErrUnknown, ErrDb), ErrDb))
	// fmt.Println(errors.Is(ErrDb, ErrDb))
	// fmt.Println(erx.Join(ErrDb, ErrUnknown))

	e := erx.W(ErrDb, "db error2")
	// fmt.Println(e)
	// fmt.Println(errors.Join(e, ErrUnknown))
	fmt.Println(erx.GetClientMsg(e))
	fmt.Println(erx.GetFullMsg(e))
	// erx.Log(ErrUnknown)
	// erx.Log(erx.Join(ErrUnknown, ErrDb))

	// a := erx.New("gg")
	// erx.Log(erx.W(ErrUnknown))

	errDb := errors.New("err db")
	errUnknow := errors.New("err unknow")
	e2 := errors.Join(errDb, errUnknow)
	fmt.Println(errors.Is(e2, errors.New("err unknow")))
}
