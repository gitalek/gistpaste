package helpers

import (
	"github.com/gitalek/gistpaste/pkg/errors"
	"net/http"
	"strconv"
)

func ValidateParamId(p string) (int, *errors.ErrBadRequest) {
	if p == "" {
		return 0, &errors.ErrBadRequest{
			StatusCode: http.StatusBadRequest,
			Message:    "id query param is not found"}
	}
	id, err := strconv.Atoi(p)
	if err != nil {
		return 0, &errors.ErrBadRequest{
			StatusCode: http.StatusBadRequest,
			Message:    "your id query param is not a number [should be positive integer]",
		}
	}
	if id < 1 {
		return 0, &errors.ErrBadRequest{
			StatusCode: http.StatusBadRequest,
			Message:    "your id query param is less than or equal to zero [should be more than or equal to one]",
		}
	}
	return id, nil
}
