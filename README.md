# 🌟 erx — Minimal-Effort, Rich-Context Error Handling for Go

<p align="center">
  <a href="https://github.com/skyrocket-qy/erx/actions/workflows/ci.yml"><img src="https://github.com/skyrocket-qy/erx/actions/workflows/ci.yml/badge.svg" alt="Build Status"></a>
  <a href="#"><img src="https://img.shields.io/badge/coverage-95.3%25-brightgreen" alt="Coverage"></a>
</p>

`erx` provides a simple way to create and wrap errors, enriching them with call stacks and custom error codes without cluttering your code. It's designed for easy integration and clear, structured error reporting.

---

## ✨ Features

- **🔢 Custom Error Codes:** Define and attach unique codes to your errors for easy identification and client-side handling.
- **🧠 Automatic Call Stacks:** Automatically capture a full call stack when an error is created or wrapped, without any extra effort.
- **🔄 Effortless Error Wrapping:** Wrap existing errors from any source (`os`, `database/sql`, etc.) to add context and a call stack.
- **🤝 Standard Library Compatibility:** Works seamlessly with the standard `errors` package, including `errors.Is()` and `errors.As()`.
- **🔧 Customizable Error Mapping:** Provide a custom mapping function to convert third-party errors into specific `erx` error codes automatically.
- **💬 Rich Context:** Add descriptive messages at each level of the call stack to build a clear narrative of the error's path.

---

## 📦 Installation

To install `erx`, use `go get`:
```sh
go get github.com/skyrocket-qy/erx
```

---

## 🚀 Usage

First, define your custom error codes. It's a good practice to keep them in a central `errors.go` file.

```go
package app

import "github.com/skyrocket-qy/erx"

const (
	ErrUnknown      erx.CodeImp = "500.0000" // Internal Server Error
	ErrInvalidInput erx.CodeImp = "400.0001" // Bad Request
	ErrNotFound     erx.CodeImp = "404.0001" // Not Found
)
```

### Creating New Structured Errors

Use `erx.New()` or `erx.Newf()` when you are the originator of the error. This creates a new `*erx.CtxErr` with a call stack and your custom error code.

```go
import (
    "fmt"
    "github.com/skyrocket-qy/erx"
)

func findUser(id int) (*User, error) {
    if id <= 0 {
        // Create a new error with a code and a formatted message
        return nil, erx.Newf(ErrInvalidInput, "user ID must be positive, got %d", id)
    }
    // ...
}
```

### Wrapping Existing Errors

When you receive an error from another function (e.g., a database driver or a standard library function), use `erx.W()` to wrap it. This preserves the original error while enriching it with a call stack.

```go
import (
    "database/sql"
    "github.com/skyrocket-qy/erx"
)

func getUserFromDB(id int) (*User, error) {
    user := &User{}
    err := db.QueryRow("SELECT id, name FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name)
    if err != nil {
        if err == sql.ErrNoRows {
            // Wrap the original error and assign a more specific code
            return nil, erx.W(err, "user not found in database").SetCode(ErrNotFound)
        }
        // Wrap the original error, letting the default mapper handle the code
        return nil, erx.W(err, "database query failed")
    }
    return user, nil
}
```

### Inspecting Structured Errors

An `*erx.CtxErr` contains valuable debugging information. You can access the error code, the call stack, and any messages attached along the way.

```go
import (
    "errors"
    "fmt"
    "github.com/skyrocket-qy/erx"
)

func main() {
    _, err := findUser(0)
    if err != nil {
        var ctxErr *erx.CtxErr
        // Use errors.As to check if the error is an erx error
        if errors.As(err, &ctxErr) {
            fmt.Printf("Error Code: %s\n", ctxErr.Code.Str())
            fmt.Println("Call Stack:")
            for _, frame := range ctxErr.CallerInfos {
                fmt.Printf("  - %s:%d\n    %s\n", frame.File, frame.Line, frame.Function)
                if frame.Msg != "" {
                    fmt.Printf("    Message: %s\n", frame.Msg)
                }
            }
        }
    }
}
```

### Working with Standard `errors`

`erx` is fully compatible with the standard `errors` package, so you can use `errors.Is()` and `errors.As()` as you normally would.

#### `errors.Is()`

Use `errors.Is()` to check if an error in the wrap chain matches a specific error instance.

```go
import (
    "database/sql"
    "errors"
)

func main() {
    _, err := getUserFromDB(123) // Assuming this user doesn't exist
    if err != nil {
        // Check if the underlying cause was sql.ErrNoRows
        if errors.Is(err, sql.ErrNoRows) {
            fmt.Println("User not found (checked with errors.Is)")
        }
    }
}
```

#### `errors.As()`

Use `errors.As()` to access the underlying `*erx.CtxErr` to inspect its contents. This is the preferred way to get access to the structured error.

```go
import (
    "errors"
    "github.com/skyrocket-qy/erx"
)

func main() {
    err := someFunctionThatReturnsAnError()
    var ctxErr *erx.CtxErr
    if errors.As(err, &ctxErr) {
        // Now you can work with the structured error
        if ctxErr.Code == ErrNotFound {
            // ...
        }
    }
}
```

---

## 🧠 Advanced Usage

### Custom Error Code Mapping

`erx` can automatically assign a custom error code when wrapping an error from a third-party library. To enable this, you need to overwrite the global `erx.ErrToCode` function with your own logic.

This is useful for mapping specific database errors, I/O errors, etc., to your application's error codes without writing `if/else` blocks everywhere.

```go
import (
    "errors"
    "gorm.io/gorm" // Example using GORM
    "github.com/skyrocket-qy/erx"
)

func init() {
    // Overwrite the default mapping function.
    // This should be done once when your application starts.
    erx.ErrToCode = func(err error) erx.Code {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrNotFound
        }
        // Add more mappings here...

        // Fall back to the default unknown error code.
        return erx.ErrUnknown
    }
}

// Now, when you wrap gorm.ErrRecordNotFound, it will automatically get the ErrNotFound code.
func someDatabaseCall() error {
    err := gorm.ErrRecordNotFound // Simulate a GORM error
    // erx.W will call our custom ErrToCode function internally.
    return erx.W(err, "failed to find record")
}
```

---

## 📖 API Reference

- `erx.New(code Code, msgs ...string) *CtxErr`
  Creates a new structured error with a code, a message, and a call stack.

- `erx.Newf(code Code, format string, args ...any) *CtxErr`
  Creates a new structured error with a code, a formatted message, and a call stack.

- `erx.W(err error, msgs ...string) *CtxErr`
  Wraps an existing error, adding a call stack and an optional message. If the error is not already an `*erx.CtxErr`, it will be assigned a code by the `erx.ErrToCode` function.
---

## ⚖️ License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
