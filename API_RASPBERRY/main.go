// @title API de Atracciones y Visitas
// @version 1.0
// @description API para registrar y consultar atracciones y visitas, integradas con RabbitMQ.
// @termsOfService http://swagger.io/terms/
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
package main

import (
	"api1/core/database"
	_ "api1/docs"
	"api1/src/views"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()
	views.RegisterRoutes(r)

	port := "8080"
	fmt.Printf("\n🚀 Servidor iniciado en http://localhost:%s\n", port)
	fmt.Printf("📘 Documentación Swagger disponible en: http://localhost:%s/swagger/index.html\n\n", port)

	r.Run(":" + port)
}

