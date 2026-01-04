package handler

import (
	"best-pattern/internal/middleware"
	"best-pattern/internal/response"
	"best-pattern/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BookHandler struct {
	bookService service.BookServiceInteface
}

type BookHandlerInterface interface {
	Routes() http.Handler
}

func NewBookHandler(bookService service.BookServiceInteface) BookHandlerInterface {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) Routes() http.Handler {
	r := chi.NewRouter()

	getResourceID := func(r *http.Request) string {
		return chi.URLParam(r, "id")
	}

	auditMiddleware := func(action, resourceType string) func(http.Handler) http.Handler {
		return middleware.AuditMiddleware(action, resourceType, getResourceID)
	}

	r.With(auditMiddleware("/list-book", "book")).Get("/", h.ListBook)

	return r

}

func (h *BookHandler) ListBook(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sortBy := r.URL.Query().Get("sort_by")
	orderBy := r.URL.Query().Get("order_by")
	search := r.URL.Query().Get("search")

	filters := make(map[string]string)
	filters["name"] = r.URL.Query().Get("name")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 10
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 1
	}

	books, total, err := h.bookService.ListBook(ctx, filters, search, page, limit, sortBy, orderBy)
	if err != nil {
		response := response.NewErrorResponse(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Calculate pagination
	pagination := response.Pagination{
		CurrentPage: page,
		From:        (page-1)*limit + 1,
		To:          (page-1)*limit + len(books),
		Pages:       (total + limit - 1) / limit,
		Total:       total,
	}

	// Create success response with pagination
	successResponse := response.NewSuccessResponseWithPagination(books, pagination)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(successResponse)
}
