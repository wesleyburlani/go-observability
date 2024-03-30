package utils

import (
	"errors"
	"net/http"

	pkg_errors "github.com/wesleyburlani/go-observability/pkg/errors"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	statusCode := resolveStatusCode(err)
	SendJsonResponse(w, r, statusCode, ErrorResponse{Error: err.Error()})
}

func resolveStatusCode(err error) int {
	if errors.Is(err, pkg_errors.ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, pkg_errors.ErrUnauthorized) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, pkg_errors.ErrValidation) {
		return http.StatusBadRequest
	}
	if errors.Is(err, pkg_errors.ErrConflict) {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}
