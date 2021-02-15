package helpers

import (
	"github.com/gitalek/gistpaste/pkg/errors"
	"net/http"
	"strconv"
)

func ValidateParamId(p string) *errors.ErrBadRequest {
	if p == "" {
		return &errors.ErrBadRequest{
			StatusCode: http.StatusBadRequest,
			Message:    "id query param is not found"}
	}
	id, err := strconv.Atoi(p)
	if err != nil {
		return &errors.ErrBadRequest{
			StatusCode: http.StatusBadRequest,
			Message:    "your id query param is not a number [should be positive integer]",
		}
	}
	if id < 1 {
		return &errors.ErrBadRequest{
			StatusCode: http.StatusBadRequest,
			Message:    "your id query param is less than or equal to zero [should be more than or equal to one]",
		}
	}
	return nil
}
