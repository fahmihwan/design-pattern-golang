package request

import (
	"best-pattern/internal/model"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

type BookRequest struct {
	Title       string  `json:"name" validate:"required"`
	Author      string  `json:"author" validate:"required"`
	Description *string `json:"description,omitempty"`
}

func (r *BookRequest) parse(req *multipart.Form) {
	var values = req.Value
	r.Title = getStringFrom(values["title"])
	r.Author = getStringFrom(values["author"])
	desc := getStringFrom(values["description"])

	if desc != "" {
		r.Description = &desc
	} else {
		r.Description = nil
	}
}

func (r *BookRequest) validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (r *BookRequest) ToBook() *model.Book {
	return &model.Book{
		Title:       r.Title,
		Author:      r.Author,
		Description: r.Description,
	}
}
