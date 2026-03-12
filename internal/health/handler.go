package health

import (
	"net/http"

	"go-api/internal/gen/openapi"
	"go-api/internal/httpx"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetHealth(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, openapi.HealthResponse{
		Status: "ok",
	})
}
