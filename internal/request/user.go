package request

import (
	"best-pattern/internal/model"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UserRegisterRequest struct {
	Name     string  `json:"name" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password *string `json:"password,omitempty"`
}

func (r *UserRegisterRequest) parse(req *multipart.Form) {
	values := req.Value

	r.Name = strings.TrimSpace(getStringFrom(values["name"]))
	r.Email = strings.ToLower(strings.TrimSpace(getStringFrom(values["email"])))
	password := strings.TrimSpace(getStringFrom(values["password"]))

	if password != "" {
		r.Password = &password
	} else {
		r.Password = nil
	}

}

func (r *UserRegisterRequest) validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// ToUser: sengaja TIDAK set password hash di sini.
// Hashing password lakukan di service.
func (r *UserRegisterRequest) ToUser() *model.User {
	return &model.User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

// LOGIN
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r *UserLoginRequest) parse(req *multipart.Form) {
	values := req.Value

	r.Email = strings.ToLower(strings.TrimSpace(getStringFrom(values["email"])))
	r.Password = strings.TrimSpace(getStringFrom(values["password"]))
}

func (r *UserLoginRequest) validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
