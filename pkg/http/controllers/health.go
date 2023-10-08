package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

type Health struct{}

func NewHealth() *Health {
	return &Health{}
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := HealthGetResponse{
		Status: "healthy",
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
