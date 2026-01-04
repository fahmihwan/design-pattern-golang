package request

import (
	"best-pattern/internal/response"
	"errors"
	"mime/multipart"
	"net/http"
	"net/url"
)

type PaginationRequest[T any] struct {
	Page    int    `json:"page" default:"1"`
	Limit   int    `json:"limit" default:"10"`
	SortBy  string `json:"sort_by"`
	OrderBy string `json:"order_by" default:"DESC"`
	Search  string `json:"search"`
	Filter  T      `json:"filter"`
}

func (r *PaginationRequest[T]) GetOffset() int {
	return (r.Page - 1) * r.Limit
}

func (r *PaginationRequest[T]) ToPagination(length, total int) response.Pagination {
	return response.Pagination{
		CurrentPage: r.Page,
		From:        (r.Page-1)*r.Limit + 1,
		To:          (r.Page-1)*r.Limit + length,
		Pages:       (total + r.Limit - 1) / r.Limit,
		Total:       total,
	}
}

type RequestParse interface {
	validate() error
	parse(req *multipart.Form)
}

func ParseForm(r *http.Request, data RequestParse) (err error) {
	if err = r.ParseMultipartForm(32 << 20); err != nil {
		return
	}

	data.parse(r.MultipartForm)
	if data == nil {
		return errors.New("invalid request type")
	}

	if err = data.validate(); err != nil {
		return
	}

	return nil
}

type RequestParser interface {
	ParseForm(*multipart.Form)
}

// Helper function for handlers
func ParseMultipartRequest(r *http.Request, req RequestParser) error {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return err
	}
	req.ParseForm(r.MultipartForm)
	return nil
}

type RequestFilter interface {
	ParseFilter(query url.Values)
}

func ParsePagination[T RequestFilter](r *http.Request, filter T) (req PaginationRequest[T]) {
	var query = r.URL.Query()
	req.Page = getIntFromQuery(query.Get("page"), 1)
	req.Limit = getIntFromQuery(query.Get("limit"), 10)
	req.SortBy = query.Get("sort_by")
	req.OrderBy = getStringOrDefault(query.Get("order_by"), "DESC")
	req.Search = query.Get("search")

	filter.ParseFilter(query)
	req.Filter = filter

	return
}

type Filter struct{}

func (*Filter) ParseFilter(query url.Values) {}
