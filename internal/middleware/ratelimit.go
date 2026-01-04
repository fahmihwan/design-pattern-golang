package middleware

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

func RateLimitingMiddleware(next http.Handler) http.Handler {
	limiter := tollbooth.NewLimiter(100, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	return tollbooth.LimitHandler(limiter, next)
}
