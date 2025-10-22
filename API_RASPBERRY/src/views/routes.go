package views

import (
	_ "api1/docs"
	"api1/src/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine) {
	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Endpoints principales
	r.POST("/visitas", controllers.CreateVisitas)
	r.POST("/atracciones", controllers.CreateAtracciones)
	r.GET("/visitas/dia", controllers.GetVisitasFromDate)
	r.GET("/atracciones/dia", controllers.GetAtraccionesFromDate)
}
