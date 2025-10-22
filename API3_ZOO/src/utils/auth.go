package utils

import (
	"net/http"
	"strings"
)

func RequireRole(allowedRoles ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.URL.Query().Get("token")
			if tokenStr == "" {
				http.Error(w, "Token requerido en la URL", http.StatusUnauthorized)
				return
			}

			claims, err := ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}

			for _, role := range allowedRoles {
				if strings.EqualFold(claims.Role, role) {
					next(w, r)
					return
				}
			}

			http.Error(w, "Acceso no autorizado: rol insuficiente", http.StatusForbidden)
		}
	}
}
