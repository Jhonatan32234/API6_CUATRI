package controllers

import (
	"api1/src/entities"
	"api1/src/models"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAtraccionesFromDate godoc
// @Summary Obtener atracciones desde una fecha
// @Description Devuelve atracciones a partir de una fecha, y las envía a RabbitMQ
// @Tags atracciones
// @Produce json
// @Param fecha query string true "Fecha desde la cual obtener atracciones (YYYY-MM-DD)"
// @Success 200 {array} entities.Atraccion
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /atracciones/dia [get]
func GetAtraccionesFromDate(c *gin.Context) {
	fecha := c.Query("fecha")
	if fecha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere parámetro 'fecha'"})
		return
	}

	atracciones, err := models.GetAtraccionesFromDate(fecha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener atracciones"})
		return
	}

	c.JSON(http.StatusOK, atracciones)
}

// CreateAtracciones godoc
// @Summary Crear una o varias atracciones
// @Description Crea atracciones en la base de datos y las envía a RabbitMQ si es posible
// @Tags atracciones
// @Accept json
// @Produce json
// @Param body body []entities.Atraccion true "Arreglo de atracciones o una sola atracción"
// @Success 201 {array} entities.Atraccion
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /atracciones [post]
func CreateAtracciones(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No se pudo leer el cuerpo"})
		return
	}

	var unaAtraccion entities.Atraccion
	if err := json.Unmarshal(body, &unaAtraccion); err == nil {
		result, err := models.SaveAtracciones([]entities.Atraccion{unaAtraccion})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar atracción"})
			return
		}
		c.JSON(http.StatusCreated, result)
		return
	}

	var variasAtracciones []entities.Atraccion
	if err := json.Unmarshal(body, &variasAtracciones); err == nil {
		result, err := models.SaveAtracciones(variasAtracciones)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar atracciones"})
			return
		}
		c.JSON(http.StatusCreated, result)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
}
