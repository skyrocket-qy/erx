package erx

// func In(err error, targets []error) bool {
// 	for _, t := range targets {
// 		if errors.Is(err, t) {
// 			return true
// 		}
// 	}
// 	return false
// }

// Recursively unwraps the error to get the root cause.
// only support single chain of errors
// func RUnwrap(err error) error {
// 	for {
// 		unwrapped := errors.Unwrap(err)
// 		if unwrapped == nil {
// 			return err
// 		}
// 		err = unwrapped
// 	}
// }

// From top-down error, check if it can map to non-500 http code
// return it, if not found, return 500
// not support single chain of errors
// func ToHttpCode(err error, mapping func(err error) int) int {
// 	if errCtx, ok := err.(*ContextError); ok {
// 		err = errCtx.Unwrap()
// 	}

// 	if jErr, ok := err.(interface{ Unwrap() []error }); ok {
// 		joinErrs := jErr.Unwrap()
// 		for _, e := range joinErrs {
// 			if code := mapping(e); code != http.StatusInternalServerError {
// 				return code
// 			}
// 		}
// 	} else {
// 		return mapping(err)
// 	}

// 	return http.StatusInternalServerError
// }
