package main

import (
	"errors"
	"fmt"

	"github.com/skyrocketOoO/erx/erx"
)

var (
	ErrUnknown = erx.New("unknown error")
	ErrDb      = erx.New("db error")
)

var OgErrDb = errors.New("db error")

func main() {
	fmt.Println(erx.Is(fmt.Errorf("%w", ErrDb), ErrDb))
	fmt.Println(erx.Is(erx.Wrap("db error", ErrUnknown), ErrUnknown))
	fmt.Println(erx.Is(erx.Join(ErrUnknown, ErrDb), ErrDb))
	fmt.Println(erx.Is(ErrDb, ErrDb))
}
