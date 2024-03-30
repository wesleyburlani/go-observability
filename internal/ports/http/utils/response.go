package utils

import (
	"net/http"

	"github.com/go-chi/render"
)

func SendJsonResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	render.JSON(w, r, data)
}
