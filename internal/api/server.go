package api

import (
	"net/http"

	"go-api/internal/gen/openapi"
	"go-api/internal/health"
	"go-api/internal/user"
)

type Handler struct {
	healthHandler *health.Handler
	userHandler   *user.Handler
}

func NewHandler(healthHandler *health.Handler, userHandler *user.Handler) *Handler {
	return &Handler{
		healthHandler: healthHandler,
		userHandler:   userHandler,
	}
}

func (h *Handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	h.healthHandler.GetHealth(w, r)
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request, params openapi.ListUsersParams) {
	h.userHandler.ListUsers(w, r, params)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.userHandler.CreateUser(w, r)
}

func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request, id int64) {
	h.userHandler.GetUserByID(w, r, id)
}

func (h *Handler) UpdateUserById(w http.ResponseWriter, r *http.Request, id int64) {
	h.userHandler.UpdateUserByID(w, r, id)
}

func (h *Handler) DeleteUserById(w http.ResponseWriter, r *http.Request, id int64) {
	h.userHandler.DeleteUserByID(w, r, id)
}
