package middleware

import (
	"best-pattern/internal/util"
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const (
	CtxUserID ctxKey = "user_id"
	CtxEmail  ctxKey = "email"
	// CtxRole   ctxKey = "role"
)

func AuthJWT(jwtm *util.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := r.Header.Get("Authorization")

			if h == "" || !strings.HasPrefix(h, "Bearer ") {
				http.Error(w, "missing authorization token", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(h, "Bearer ")
			claims, err := jwtm.Verify(tokenString)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxEmail, claims.Email)
			// ctx = context.WithValue(ctx, CtxRole, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

// func AuthJWT(jwtm *util.JWTManager) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			auth := strings.TrimSpace(r.Header.Get("Authorization"))
// 			if auth == "" {
// 				http.Error(w, "unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			// split: "Bearer <token>"
// 			parts := strings.SplitN(auth, " ", 2)
// 			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
// 				http.Error(w, "unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			tokenString := strings.TrimSpace(parts[1])

// 			claims, err := jwtm.Verify(tokenString)
// 			if err != nil {
// 				http.Error(w, "unauthorized", http.StatusUnauthorized)
// 				return
// 			}

// 			ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
// 			ctx = context.WithValue(ctx, CtxEmail, claims.Email)
// 			// ctx = context.WithValue(ctx, CtxRole, claims.Role)

// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }
