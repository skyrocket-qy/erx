# 🌟 erx — Structured Error Handling for Go

`erx` is a lightweight Go library for error handling with structured codes, automatic call stack capture with additional context, and clean separation between client-facing and internal error information.

---

## ✨ Features

- 🔢 **Custom Error Codes** (`Coder`) — safely expose to clients with i18n handling
- 🧩**Extra Context** - support key-value paris or single string information on each call
- 🧠 **Call Stack Capture** — for better debugging
- 🔄 **Error Wrapping Helpers** — contextualize errors easily
- 🔐 **Client vs Internal Separation** — clean boundary of what to expose

---

## 🚀 Quick Start

Define interanl error code

```go
const (
  ErrUnknown erx.CoderImp = "500.0000"
  ErrNotFound erx.CoderImp = "404.0000"
)
```

Overwrite the default ErrToCode function

```go
// use to convert 3rd party error to erx.Coder
erx.ErrToCode = func(err error) erx.Coder {
  if errors.As(err, gorm.ErrRecordNotFound) {
    return ErrNotFound
  }

  return ErrUnknown
}
```

New or wrap error or add something

```go
if err != nil{
  return erx.W(err, "another context") // or return err if this is not first error
}

if !ok {
  return erx.New(erx.ErrNotFound, fmt.Sprintf("key: %s", key))
}

// for non-defined mapping 3rd party error
if err := thirdPartyFunc(); err != nil {
  return erx.WCode(err, Erx.ErrNotFound, "another context")
}
```
