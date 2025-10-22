package controllers

import (
	"api2/src/entities"
	"api2/src/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateVisitaGeneral godoc
// @Summary Crear una nueva visita
// @Tags visitasGeneral
// @Accept json
// @Produce json
// @Security JWTQuery
// @Param visita body entities.VisitaGeneral true "Datos de la visita"
// @Success 201 {object} map[string]string
// @Failure 401 {object} map[string]string "Token inválido o no autorizado"
// @Failure 400 {object} map[string]string
// @Router /visitasgeneral [post]
func CreateVisitaGeneral(c *gin.Context) {
	var visita entities.VisitaGeneral

	if err := c.ShouldBindJSON(&visita); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Asignar fecha actual (solo la parte de la fecha, sin hora)
	visita.Fecha = time.Now().Format("2006-01-02")

	if err := models.CreateVisitaGeneral(visita); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Visita creada"})
}

// GetAllVisitasGeneral godoc
// @Summary Obtener todas las visitas registradas
// @Tags visitasGeneral
// @Produce json
// @Security JWTQuery
// @Success 200 {array} entities.VisitaGeneral
// @Failure 401 {object} map[string]string "Token inválido o no autorizado"
// @Failure 500 {object} map[string]string
// @Router /visitasgeneral [get]
func GetAllVisitasGeneral(c *gin.Context) {
	data, err := models.GetAllVisitasGeneral()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetVisitaGeneralByFecha godoc
// @Summary Obtener una visita por fecha
// @Tags visitasGeneral
// @Produce json
// @Security JWTQuery
// @Param fecha path string true "Fecha de la visita en formato YYYY-MM-DD"
// @Success 200 {object} entities.VisitaGeneral
// @Failure 401 {object} map[string]string "Token inválido o no autorizado"
// @Failure 404 {object} map[string]string
// @Router /visitasgeneral/{fecha} [get]
func GetVisitaGeneralByFecha(c *gin.Context) {
	fecha := c.Param("fecha")
	var visita entities.VisitaGeneral
	err := models.GetVisitaGeneralByFecha(fecha, &visita)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Visita no encontrada"})
		return
	}
	c.JSON(http.StatusOK, visita)
}


// UpdateVisitaGeneral godoc
// @Summary Actualizar una visita por fecha
// @Tags visitasGeneral
// @Accept json
// @Produce json
// @Security JWTQuery
// @Param fecha path string true "Fecha de la visita en formato YYYY-MM-DD"
// @Param visita body entities.VisitaGeneral true "Nuevos datos de la visita"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string "Token invál
// @Failure 400 {object} map[string]string
// @Router /visitasgeneral/{fecha} [put]
func UpdateVisitaGeneral(c *gin.Context) {
	fecha := c.Param("fecha") // fecha en formato "2025-07-23"
	var updated entities.VisitaGeneral
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateVisitaGeneralPorFecha(fecha, updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Visita actualizada"})
}


// DeleteVisitaGeneralPorFecha godoc
// @Summary Eliminar una visita por fecha
// @Tags visitasGeneral
// @Produce json
// @Security JWTQuery
// @Param fecha path string true "Fecha de la visita en formato YYYY-MM-DD"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string "Token inválido o no autorizado"
// @Failure 500 {object} map[string]string
// @Router /visitasgeneral/{fecha} [delete]
func DeleteVisitaGeneralPorFecha(c *gin.Context) {
	fecha := c.Param("fecha")
	if err := models.DeleteVisitaGeneralPorFecha(fecha); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Visita eliminada"})
}

