package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	http_utils "github.com/wesleyburlani/go-rest/internal/transport/http/utils"
	"github.com/wesleyburlani/go-rest/internal/users"
	"github.com/wesleyburlani/go-rest/pkg/logger"
)

type Users struct {
	svc    *users.Service
	logger *logger.Logger
}

type CreateUserBody struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type UpdateUserBody struct {
	Username string `json:"username" validate:"omitempty,min=3,max=20"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8,max=20"`
}

type LoginBody struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

func NewUsers(svc *users.Service, logger *logger.Logger) *Users {
	return &Users{svc: svc, logger: logger}
}

func (c *Users) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/{id}", c.get)
	r.Post("/", c.create)
	r.Put("/{id}", c.update)
	r.Delete("/{id}", c.delete)
	r.Post("/login", c.login)
	return r
}

func (c *Users) get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	id, err := http_utils.GetInt64UrlParam(r, "id")
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}

	u, err := c.svc.Get(ctx, id)
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}
	http_utils.SendJsonResponse(w, r, http.StatusOK, u)
}

func (c *Users) create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	body, err := http_utils.ParseBody[CreateUserBody](r)
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}

	u, err := c.svc.Create(ctx, users.User{
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}
	http_utils.SendJsonResponse(w, r, http.StatusOK, u)
}

func (c *Users) update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	id, err := http_utils.GetInt64UrlParam(r, "id")
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}

	body, err := http_utils.ParseBody[UpdateUserBody](r)
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}

	u, err := c.svc.Update(ctx, users.User{
		ID:       id,
		Username: body.Username,
		Email:    body.Email,
		Password: body.Password,
	})
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}
	http_utils.SendJsonResponse(w, r, http.StatusOK, u)
}

func (c *Users) delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	id, err := http_utils.GetInt64UrlParam(r, "id")
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}

	_, err = c.svc.Delete(ctx, id)
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}
	http_utils.SendJsonResponse(w, r, http.StatusOK, nil)
}

func (c *Users) login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()

	body, err := http_utils.ParseBody[LoginBody](r)
	if err != nil {
		http_utils.HandleError(w, r, err)
		return
	}

	if err := c.svc.Login(ctx, body.Username, body.Password); err != nil {
		http_utils.HandleError(w, r, err)
		return
	}
	http_utils.SendJsonResponse(w, r, http.StatusOK, nil)
}
