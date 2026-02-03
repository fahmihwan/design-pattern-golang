package handler

import (
	"best-pattern/internal/middleware"
	"best-pattern/internal/request"
	"best-pattern/internal/response"
	"best-pattern/internal/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService service.UserServiceInterface
}
type UserHandlerInterface interface {
	Routes() http.Handler
}

func NewUserHandler(userService service.UserServiceInterface) UserHandlerInterface {
	return &UserHandler{
		userService: userService,
	}
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

func (h *UserHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req = new(request.UserRegisterRequest)
	if err := request.ParseForm(r, req); err != nil {
		middleware.HandleValidationErrors(err, w)
		return
	}

	user := req.ToUser()

	createUser, err := h.userService.Register(ctx, user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating form: %v", err), http.StatusInternalServerError)
		return
	}
	response := response.NewSuccessResponse(createUser)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
