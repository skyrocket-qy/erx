# erx

A compatible built-in error but support call stack

!!!Not support single chain error!!!

## Quick start

Define error type

```go
const
var (
  ErrDB = errors.New("db error")
  ErrUnknown = errors.New("unknown error")
)
```

New or wrap error or add something

```go
if err != nil{
  return erx.Errorf("somethin wrong: %w", err)
}

if err != nil{
  return erx.W(err)
}

if !ok {
  return erx.New("something wrong")
}
```

## Scenario

When debug, the error message sometimes hard for developer to find the occur place in code.

In general, it has belows methods:

```go
if err != nil{
  log(err.Error)
  return err
}
```

Or

```go
func Log(err error) error{
  log(err.Error)
  return err
}

if err != nil{
  return Log(err)
}
```

The Problem is:

- if A() call B() call C() call D(), when error occured, above methods will log 4 times repeated, but in fact, what we need is only one call stack from deeper error
- When error message show on client side, it is hard to find the real place occured(string search? if most place use same error?)

## Target

- When first error occured, record the call stack
- Without broken original interface but support callstack
- Call stack should not propogate to client, so distinguish client msg and original msg

## Features

- errorcode
- Call stack
- Separation of client error and log errors
- error msg wrap
