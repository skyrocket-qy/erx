package main

import (
	"errors"
	"fmt"
)

type Error struct {
	CallStack string
	Err       error
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Is(err error, target error) bool {
	var a, b error
	e1, ok := err.(Error)
	if ok {
		a = e1.Err
	} else {
		a = err
	}
	e2, ok := target.(Error)
	if ok {
		b = e2.Err
	} else {
		b = target
	}

	return errors.Is(a, b)
}

func As(err error, target any) bool {
	ce, ok := err.(Error)
	if ok {
		return errors.As(ce.Err, target)
	}

	return errors.As(err, target)
}

func Join(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	var err error
	for _, e := range errs {
		if e == nil {
			continue
		}
		e1, ok := e.(Error)
		if ok {
			err = errors.Join(err, e1.Err)
		} else {
			err = errors.Join(err, e)
		}
	}
	return err
}

func Wrap(msg string, err error) error {
	return Join(errors.New(msg), err)
}

func New(text string) Error {
	return Error{
		CallStack: getCallStack(),
		Err:       errors.New(text),
	}
}

func Errorf(format string, args ...interface{}) error {
	return New(fmt.Sprintf(format, args...))
}

var (
	ErrUnknown = New("unknown error")
	ErrDb      = New("db error")
)
