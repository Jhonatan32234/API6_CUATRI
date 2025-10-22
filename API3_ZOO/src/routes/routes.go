package routes

import (
	"api3/src/controllers"
	"api3/src/utils"
	"net/http"

	"github.com/gorilla/mux"
)

// Middleware CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Cabeceras CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // Cambia "*" por tu dominio si quieres
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Preflight (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Aplica el middleware CORS globalmente
	r.Use(corsMiddleware)

	// Define rutas después de aplicar el middleware
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/register", controllers.Register).Methods("POST")
	r.HandleFunc("/update/{id}", utils.RequireRole("admin")(controllers.UpdateUser)).Methods("PUT")
	r.HandleFunc("/delete/{id}", utils.RequireRole("admin")(controllers.DeleteUser)).Methods("DELETE")
	r.HandleFunc("/users", utils.RequireRole("admin")(controllers.GetAllUsers)).Methods("GET")

	return r
}
