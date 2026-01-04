package middleware

import (
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func AuditMiddleware(action string, resourceType string, getResourceID func(r *http.Request) string,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap the resHandlerFuncponse writer

			rw := &responseWriter{w, http.StatusOK}

			// call next handler (PAKAI rw biar status/size bisa kecapture)
			next.ServeHTTP(rw, r)

		})
	}
}
