package main

import (
	"api2/db"
	_ "api2/docs"
	"api2/src/controllers"
	"api2/src/models/services"
	"api2/src/routes"
	"api2/utils"
	"api2/websocket"
	"log"

	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Monitoreo de Visitas y Atracciones
// @version 1.0
// @description API que proporciona datos en tiempo real sobre visitas y atracciones, protegida con autenticación por token JWT PARA PODER USAR LOS METODOS SE NECESITARA UN TOKEN EL CUAL UNA VEZ OBTENIDO SE INSERTARA EN EL APARTADO DE AUTHORIZE MAS ESPECIFICO EL BOTON VERDE DE ABAJO CON UN CANDADO INSERTAR EL TOKEN EN EL APARTADO VALUE Y PULSAR EL BOTON DE AUTHORIZE UNA VEZ HECHO ESO DEBERIA DE PODER CERRAR LA VENTANA SIN PROBLEMAS Y USAR LOS METODOS DE AHORA EN ADELANTE
// @BasePath /api
// @securityDefinitions.apikey JWTQuery
// @in query
// @name token
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("⚠️  Advertencia: no se pudo cargar el archivo .env:", err)
	}

	db.Connect()
	services.StartRabbitConsumers()

	go websocket.StartBroadcaster()
	go utils.StartBroadcaster()

	r := gin.Default()

	// Documentación Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket
	r.GET("/ws", controllers.WebSocketHandler)

	// Middleware CORS
	r.Use(utils.CORSMiddleware())

	// Rutas principales
	routes.SetupRoutes(r)

	// Mensajes en consola
	log.Println("🚀 Servidor iniciado en http://localhost:8081")
	log.Println("📘 Documentación Swagger disponible en: http://localhost:8081/swagger/index.html")

	// Ejecutar servidor
	r.Run(":8081")
}
