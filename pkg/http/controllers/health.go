package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/wesleyburlani/go-observability/pkg/logger"
)

type Health struct {
	logger *logger.Logger
}

func NewHealth(logger *logger.Logger) *Health {
	return &Health{logger: logger}
}

type HealthGetResponse struct {
	Status string `json:"status"`
}

func (c *Health) Router() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", c.get)
	return r
}

func (c *Health) get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	resp := HealthGetResponse{
		Status: "healthy",
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		c.logger.
			With(err).
			Error(r.Context(), "Error happened in JSON marshal")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
