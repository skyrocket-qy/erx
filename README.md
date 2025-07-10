# 🌟 erx — Structured Error Handling for Go

`erx` is a lightweight Go library for error handling with structured codes, automatic call stack capture, and clean separation between client-facing and internal error information.

---

## ✨ Features

- 🔢 **Custom Error Codes** (`Coder`) — safely expose to clients
- 🧠 **Call Stack Capture** — for better debugging
- 🔄 **Error Wrapping Helpers** — contextualize errors easily
- 🔐 **Client vs Internal Separation** — clean boundary of what to expose

---

## 🚀 Quick Start

Define error type

```go
const (
  Unknown erx.CoderImp = "Unknown"
  NotFound erx.CoderImp = "NotFound"
)
```

Overwrite your customize ErrToCode function

```go
// use to convert 3rd party error to erx.Coder
erx.ErrToCode = func(err error) erx.Coder {
  if errors.As(err, gorm.ErrRecordNotFound) {
    return NotFound
  }

  return Unknown
}
```

New or wrap error or add something

```go
if err != nil{
  return erx.W(err, "another context")
}

if !ok {
  return erx.New(erx.NotFound, fmt.Sprintf("key: %s", key))
}
```
