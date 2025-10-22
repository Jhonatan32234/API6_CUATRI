package controllers

import (
	"api2/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
)


// GetNowVisitas godoc
// @Summary Obtener visitas de la fecha más reciente
// @Tags visitas
// @Produce json
// @Security JWTQuery
// @Success 200 {array} entities.Visitas
// @Failure 401 {object} map[string]string "Token inválido o no autorizado"
// @Failure 500 {object} map[string]string
// @Router /visitas/now [get]
func GetNowVisitas(c *gin.Context) {
	zona := c.GetString("zona")
	data, err := models.GetNowVisitas(zona)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetLastWeekVisitas godoc
// @Summary Obtener visitas de las 6 fechas más recientes
// @Tags visitas
// @Produce json
// @Security JWTQuery
// @Success 200 {array} entities.Visitas
// @Failure 401 {object} map[string]string "Token inválido o no autorizado"
// @Failure 500 {object} map[string]string
// @Router /visitas/lastweek [get]
func GetLastWeekVisitas(c *gin.Context) {
	zona := c.GetString("zona")
	data, err := models.GetLastWeekVisitas(zona)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetYesterdayVisitas godoc
// @Summary Obtener visitas de la penúltima fecha registrada
// @Tags visitas
// @Produce json
// @Security JWTQuery
// @Success 200 {array} entities.Visitas
// @Failure 401 {object} map[string]string "Token inválido o no autorizado"
// @Failure 500 {object} map[string]string
// @Router /visitas/yesterday [get]
func GetYesterdayVisitas(c *gin.Context) {
	zona := c.GetString("zona")
	data, err := models.GetYesterdayVisitas(zona)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetOjivaVisitas godoc
// @Summary Obtener ojiva de visitas (suma por hora)
// @Tags visitas
// @Produce json
// @Security JWTQuery
// @Param fecha query string false "Fecha en formato YYYY-MM-DD"
// @Failure 401 {object} map[string]string "Token inválido o no autorizado"
// @Failure 500 {object} map[string]string
// @Router /visitas/ojiva [get]
func GetOjivaVisitas(c *gin.Context) {
	zona := c.GetString("zona")
	fecha := c.Query("fecha")
	data, err := models.GetOjivaVisitas(fecha, zona)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
