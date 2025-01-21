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

	// e := erx.W(ErrDb, "db error2")
	// fmt.Println(e)
	// fmt.Println(errors.Join(e, ErrUnknown))
	// fmt.Println(erx.GetClientMsg(e))
	// fmt.Println(erx.GetFullMsg(e))
	// erx.Log(ErrUnknown)
	// erx.Log(erx.Join(ErrUnknown, ErrDb))

	// a := erx.New("gg")
	// erx.Log(erx.W(ErrUnknown))

	// errDb := errors.New("err db")
	// errUnknow := errors.New("err unknow")
	// e2 := errors.Join(errDb, errUnknow)
	// fmt.Println(errors.Is(e2, errors.New("err unknow")))

	// err1 := errors.New("error 1")
	// Create multiple errors
	err1 := errors.New("file not found")
	err2 := errors.New("network timeout")

	// Join the errors into one
	joinedErr := errors.Join(err1, err2)
	joinedErr = erx.W(joinedErr)

	// Now we want to access all the errors from the joined error.

	// Print out the errors in the chain
	// fmt.Println(fmt.Errorf("abc: %w", joinedErr))
	fmt.Println("Unwrapped errors:", erx.Cause(joinedErr))
	// for _, e := range unwrappedErrors {
	// 	fmt.Println("- " + e.Error())
	// }
}

func Unwrap(err error) []error {
	u, ok := err.(interface {
		Unwrap() []error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}
