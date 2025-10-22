package main

import (
	"api3/db"
	_ "api3/docs" // 👈 import anónimo para Swagger
	"api3/src/routes"
	"api3/src/utils"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)


// @title API de Gestión de Usuarios y Roles
// @version 1.0
// @description API que gestiona usuarios, roles y autenticación con JWT. PARA PODER USAR LOS METODOS SE NECESITARA UN TOKEN EL CUAL UNA VEZ OBTENIDO SE INSERTARA EN EL APARTADO DE AUTHORIZE MAS ESPECIFICO EL BOTON VERDE DE ABAJO CON UN CANDADO INSERTAR EL TOKEN EN EL APARTADO VALUE Y PULSAR EL BOTON DE AUTHORIZE UNA VEZ HECHO ESO DEBERIA DE PODER CERRAR LA VENTANA SIN PROBLEMAS Y USAR LOS METODOS DE AHORA EN ADELANTE
// @BasePath /
// @securityDefinitions.apikey JWTQuery
// @in query
// @name token
// @description Ingresa el token JWT como parámetro de consulta, por ejemplo: ?token=<tu_token>
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️  Advertencia: no se pudo cargar el archivo .env:", err)
	}

	db.ConnectDB()

	// Configuración Swagger
	// (Si usas swag init, docs.SwaggerInfo se genera automáticamente)
	// docs.SwaggerInfo.BasePath = "/"

	r := routes.SetupRoutes()

	// Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Middleware CORS
	handlerWithCORS := utils.CORS(r)

	log.Println("🚀 Servidor corriendo en http://localhost:8082")
	log.Println("📘 Documentación Swagger disponible en: http://localhost:8082/swagger/index.html")

	log.Fatal(http.ListenAndServe(":8082", handlerWithCORS))
}