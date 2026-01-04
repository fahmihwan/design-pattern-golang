package handler

import (
	"best-pattern/internal/middleware"
	"best-pattern/internal/response"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandlerInterface interface {
	Routes() http.Handler
}

type UserHandler struct{}

func NewUserHandler() UserHandlerInterface {
	return &UserHandler{}
}

func (h *UserHandler) Routes() http.Handler {
	r := chi.NewRouter()

	getResourceID := func(r *http.Request) string {
		return chi.URLParam(r, "id")
	}
	auditMiddleware := func(action, resourceType string) func(http.Handler) http.Handler {
		return middleware.AuditMiddleware(action, resourceType, getResourceID)
	}

	r.With(auditMiddleware("list", "user")).Get("/", h.ListUser)

	return r
}

func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	response := response.NewSuccessResponse("wkwkw")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
