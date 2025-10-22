package controllers

import (
	"api1/src/entities"
	"api1/src/models"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetVisitasFromDate godoc
// @Summary Obtener visitas desde una fecha
// @Description Devuelve visitas a partir de una fecha, y las envía a RabbitMQ
// @Tags visitas
// @Produce json
// @Param fecha query string true "Fecha desde la cual obtener visitas (YYYY-MM-DD)"
// @Success 200 {array} entities.Visitas
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /visitas/dia [get]
func GetVisitasFromDate(c *gin.Context) {
	fecha := c.Query("fecha") // Fecha en formato string (ej. "2025-06-10")
	if fecha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere parámetro 'fecha'"})
		return
	}

	visitas, err := models.GetVisitasFromDate(fecha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener visitas"})
		return
	}

	c.JSON(http.StatusOK, visitas)
}

// CreateVisitas godoc
// @Summary Crear una o varias visitas
// @Description Crea visitas en la base de datos y las envía a RabbitMQ si es posible
// @Tags visitas
// @Accept json
// @Produce json
// @Param body body []entities.Visitas true "Arreglo de visitas o una sola visita"
// @Success 201 {array} entities.Visitas
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /visitas [post]
func CreateVisitas(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el cuerpo"})
		return
	}

	var unaVisita entities.Visitas
	if err := json.Unmarshal(body, &unaVisita); err == nil {
		result, err := models.SaveVisitas([]entities.Visitas{unaVisita})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar visita"})
			return
		}
		c.JSON(http.StatusCreated, result)
		return
	}

	var variasVisitas []entities.Visitas
	if err := json.Unmarshal(body, &variasVisitas); err == nil {
		result, err := models.SaveVisitas(variasVisitas)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar visitas"})
			return
		}
		c.JSON(http.StatusCreated, result)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
}
