package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	pkg_errors "github.com/wesleyburlani/go-rest/pkg/errors"
)

func GetStringUrlParam(r *http.Request, key string) (string, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return "", fmt.Errorf("%w: %s", pkg_errors.ErrValidation, "id is required")
	}
	return param, nil
}

func GetInt64UrlParam(r *http.Request, key string) (int64, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return 0, fmt.Errorf("%w: %s", pkg_errors.ErrValidation, "id is required")
	}
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", pkg_errors.ErrValidation, "id must be a number")
	}
	return id, nil
}

func ParseBody[T any](r *http.Request) (T, error) {
	var body T
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		return body, fmt.Errorf("%w: %s", pkg_errors.ErrValidation, err.Error())
	}

	if err := validator.New().Struct(body); err != nil {
		return body, fmt.Errorf("%w: %s", pkg_errors.ErrValidation, err.Error())
	}

	return body, nil
}
