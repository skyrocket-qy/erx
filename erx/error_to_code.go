package erx

import (
	"net/http"

	"gorm.io/gorm"
)

func errToCode(err error) int {
	switch err {
	case gorm.ErrRecordNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
