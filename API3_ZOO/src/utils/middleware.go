package utils

import (
	"net/http"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cabeceras CORS básicas
		w.Header().Set("Access-Control-Allow-Origin", "*") // o tu dominio
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Si es OPTIONS, responde y termina
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continúa con la cadena de middleware
		next.ServeHTTP(w, r)
	})
}
