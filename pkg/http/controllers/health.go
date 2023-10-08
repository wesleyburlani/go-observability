package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Health struct {
	logger *slog.Logger
}

func NewHealth(logger *slog.Logger) *Health {
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
		c.logger.Error("Error happened in JSON marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
