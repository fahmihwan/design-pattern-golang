package handler

import (
	"best-pattern/internal/middleware"
	"best-pattern/internal/request"
	"best-pattern/internal/response"
	"best-pattern/internal/service"
	"encoding/json"
	"fmt"
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

	r.With(auditMiddleware("list-book", "book")).Get("/", h.ListBook)
	r.With(auditMiddleware("create-book", "book")).Post("/", h.CreateBook)
	r.With(auditMiddleware("get-book", "book")).Get("/{id}", h.GetBook)
	r.With(auditMiddleware("update-book", "book")).Put("/{id}", h.UpdateBook)
	r.With(auditMiddleware("delete-book", "book")).Delete("/{id}", h.DeleteBook)
	return r

}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req = new(request.BookRequest)
	if err := request.ParseForm(r, req); err != nil {
		middleware.HandleValidationErrors(err, w)
		return
	}

	book := req.ToBook()
	createBook, err := h.bookService.CreateBook(ctx, book)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating form: %v", err), http.StatusInternalServerError)
		return
	}
	response := response.NewSuccessResponse(createBook)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *BookHandler) ListBook(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	sortBy := r.URL.Query().Get("sort_by")
	orderBy := r.URL.Query().Get("order_by")
	search := r.URL.Query().Get("search")

	// Set up filters
	filters := make(map[string]string)
	filters["name"] = r.URL.Query().Get("name")

	page, err := strconv.Atoi(pageStr)
	// Parse page and limit
	if err != nil {
		page = 1 // Default to page 1 if parsing fails
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10 //Default to limit of 10 if parsing fails
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

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")

	book, err := h.bookService.GetBookByID(ctx, id)
	if err != nil {
		response := response.NewErrorResponse(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	successResponse := response.NewSuccessResponse(book)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(successResponse)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating a book
	ctx := r.Context()

	var req = new(request.BookRequest)
	if err := request.ParseForm(r, req); err != nil {
		middleware.HandleValidationErrors(err, w)
		return
	}
	ids := chi.URLParam(r, "id")
	book := req.ToBook()
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	book.ID = id

	updatedBook, err := h.bookService.UpdateBook(ctx, book)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating book: %v", err), http.StatusInternalServerError)
		return
	}

	response := response.NewSuccessResponse(updatedBook)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	id := chi.URLParam(r, "id")

	err := h.bookService.DeleteBook(ctx, id)
	if err != nil {
		response := response.NewErrorResponse(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	successResponse := response.NewSuccessResponse("Book deleted successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(successResponse)
}
