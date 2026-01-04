package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Routes() http.Handler
}

type HandlerInteface struct {
	UserHandler UserHandlerInterface
	BookHandler BookHandlerInterface
	// handler lainnya banyak nanti
}

func NewRouter(handler *HandlerInteface) *chi.Mux {
	r := chi.NewRouter()

	// isinya banyak routing nanti rencananya
	r.Route("/user", func(r chi.Router) {
		r.Mount("/", handler.UserHandler.Routes())
	})

	// error di book nil pointer
	r.Route("/book", func(r chi.Router) {
		r.Mount("/", handler.BookHandler.Routes())
	})

	return r
}
