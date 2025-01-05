package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println(Is(Wrap("db error", ErrUnknown), ErrUnknown))
	fmt.Println(Is(Join(ErrUnknown, ErrDb), ErrUnknown))
	fmt.Println(Is(ErrDb, ErrDb))
}

func A() error {
	return B()
}

func B() error {
	return C()
}

func C() error {
	return errors.New(getCallStack())
}
