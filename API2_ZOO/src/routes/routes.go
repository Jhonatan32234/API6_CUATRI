package routes

import (
	"api2/src/controllers"
	"api2/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Visitas
	api.GET("/visitas/now", utils.JWTQueryMiddleware("admin", "user"), controllers.GetNowVisitas)
	api.GET("/visitas/lastweek", utils.JWTQueryMiddleware("admin", "user"), controllers.GetLastWeekVisitas)
	api.GET("/visitas/yesterday", utils.JWTQueryMiddleware("admin", "user"), controllers.GetYesterdayVisitas)
	api.GET("/visitas/ojiva", utils.JWTQueryMiddleware("admin", "user"), controllers.GetOjivaVisitas)

	// Atracciones
	api.GET("/atraccion/now", utils.JWTQueryMiddleware("admin", "user"), controllers.GetNowAtraccion)
	api.GET("/atraccion/lastweek", utils.JWTQueryMiddleware("admin", "user"), controllers.GetLastWeekAtraccion)
	api.GET("/atraccion/yesterday", utils.JWTQueryMiddleware("admin", "user"), controllers.GetYesterdayAtraccion)
	api.GET("/atraccion/ojiva", utils.JWTQueryMiddleware("admin", "user"), controllers.GetOjivaAtraccion)

	// Visitas General
	api.GET("/visitasgeneral", utils.JWTQueryMiddleware("admin","user"), controllers.GetAllVisitasGeneral)
	api.GET("/visitasgeneral/:fecha", utils.JWTQueryMiddleware("admin","user"), controllers.GetVisitaGeneralByFecha)
	api.POST("/visitasgeneral", utils.JWTQueryMiddleware("admin","user"), controllers.CreateVisitaGeneral)
	api.PUT("/visitasgeneral/:fecha", utils.JWTQueryMiddleware("admin","user"), controllers.UpdateVisitaGeneral)
	api.DELETE("/visitasgeneral/:fecha", utils.JWTQueryMiddleware("admin"), controllers.DeleteVisitaGeneralPorFecha)
}


