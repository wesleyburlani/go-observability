package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/wesleyburlani/go-rest/internal/users"
	"github.com/wesleyburlani/go-rest/pkg/logger"
)

type Users struct {
	svc    *users.Service
	logger *logger.Logger
}

func NewUsers(svc *users.Service, logger *logger.Logger) *Users {
	return &Users{svc: svc, logger: logger}
}

type HealthGetResponse struct {
	Status string `json:"status"`
}

func (c *Users) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", c.get)
	r.Post("/", c.create)
	r.Put("/{id}", c.update)
	r.Delete("/{id}", c.delete)
	return r
}

func (c *Users) get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	strId := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "id must be a number"})
		return
	}

	u, err := c.svc.Get(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, u)
}

type CreateUserBody struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func (c *Users) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	var body CreateUserBody
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	u, err := c.svc.Create(ctx, users.User{
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, u)
}

type UpdateUserBody struct {
	Username string `json:"username" validate:"omitempty,min=3,max=20"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8,max=20"`
}

func (c *Users) update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	strId := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "id must be a number"})
		return
	}

	var body UpdateUserBody
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	u, err := c.svc.Update(ctx, users.User{
		ID:       id,
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, u)
}

func (c *Users) delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	strId := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "id must be a number"})
		return
	}

	_, err = c.svc.Delete(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "deleted"})
}
