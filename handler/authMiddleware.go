package handler

import (
	"fmt"
	"net/http"
	"restapi/auth"
	"strings"
)

func (h *Handler) AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		_, err := auth.ValidateToken(parts[1])
		if err != nil {
			if err == auth.ErrInvalidToken {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			http.Error(w, fmt.Sprintf("Failed to validate token: %v", err), http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
