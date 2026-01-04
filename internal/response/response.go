package response

// Response represents the general response structure
type Response struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
}

// Pagination represents the pagination metadata
type Pagination struct {
	CurrentPage int `json:"current_page"`
	From        int `json:"from"`
	To          int `json:"to"`
	Pages       int `json:"pages"`
	Total       int `json:"total"`
}

func NewSuccessResponse(data interface{}) Response {
	return Response{
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(errors interface{}) Response {
	return Response{
		Message: "Error",
		Errors:  errors,
	}
}

// NewSuccessResponseWithPagination creates a new success response with pagination
func NewSuccessResponseWithPagination(data interface{}, pagination Pagination) Response {
	return Response{
		Message:    "Success",
		Data:       data,
		Pagination: pagination,
		Meta:       struct{}{},
	}
}
