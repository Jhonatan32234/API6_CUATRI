package models

import (
	"api1/core/database"
	"api1/core/rabbitmq"
	"api1/src/entities"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)



func GetAtraccionesFromDate(fecha string) ([]entities.Atraccion, error) {
	var atracciones []entities.Atraccion

	err := database.DB.Where("fecha >= ?", fecha).Find(&atracciones).Error
	if err != nil {
		return nil, err
	}

	if len(atracciones) == 0 {
		log.Println("No hay atracciones para enviar desde la fecha:", fecha)
		return atracciones, nil
	}

	if rabbitmq.PublishToTopic(atracciones, "atracciones_topic", "atraccion.data") {
		err = database.DB.Model(&entities.Atraccion{}).
			Where("fecha >= ? AND enviado = ?", fecha, false).
			Update("enviado", true).Error
		if err != nil {
			log.Println("Error actualizando atracciones como enviadas:", err)
		}
	}

	return atracciones, nil
}

func validateAtraccion(atraccion entities.Atraccion) error {
	// Validar que Tiempo sea mayor a 0
	if atraccion.Tiempo <= 0 {
		return fmt.Errorf("el campo 'tiempo' debe ser mayor a 0")
	}

	// Validar que no haya campos vacíos o nulos
	if atraccion.Nombre == "" || strings.TrimSpace(atraccion.Nombre) == "" {
		return fmt.Errorf("el campo 'nombre' es requerido")
	}

	if atraccion.Hora == "" || strings.TrimSpace(atraccion.Hora) == "" {
		return fmt.Errorf("el campo 'hora' es requerido")
	}

	if atraccion.Fecha == "" || strings.TrimSpace(atraccion.Fecha) == "" {
		return fmt.Errorf("el campo 'fecha' es requerido")
	}

	if atraccion.Zona == "" || strings.TrimSpace(atraccion.Zona) == "" {
		return fmt.Errorf("el campo 'zona' es requerido")
	}

	return nil
}

func SaveAtracciones(input []entities.Atraccion) ([]entities.Atraccion, error) {
	var guardadas []entities.Atraccion
	var errores []string

	for _, item := range input {
		// Validar cada atracción antes de guardar
		if err := validateAtraccion(item); err != nil {
			errores = append(errores, fmt.Sprintf("Atracción inválida: %v", err))
			log.Printf("❌ Validación fallida para atracción: %v", err)
			continue
		}

		item.Enviado = false
		if err := database.DB.Create(&item).Error; err != nil {
			log.Println("❌ Error al guardar atracción:", err)
			saveAtraccionToFile(item)
		} else {
			guardadas = append(guardadas, item)
		}
	}

	// Si hay errores de validación, retornarlos
	if len(errores) > 0 && len(guardadas) == 0 {
		return nil, fmt.Errorf("errores de validación: %s", strings.Join(errores, "; "))
	}

	if len(guardadas) == 0 {
		log.Println("⚠️ Ninguna atracción fue guardada. No se enviará al broker.")
		return nil, nil
	}

	var toSend []entities.Atraccion
	database.DB.Where("enviado = ?", false).Find(&toSend)

	if len(toSend) > 0 && rabbitmq.PublishToTopic(toSend, "atracciones_topic", "atraccion.data") {
		database.DB.Model(&entities.Atraccion{}).Where("enviado = ?", false).Update("enviado", true)

		// 🔁 Enviar por zona solo el ID
		for _, item := range toSend {
			rabbitmq.PublishIDToZoneTopic("atracciones_topic", item.Zona, item.Id, "atracciones")
		}
	}

	return toSend, nil
}

func saveAtraccionToFile(data entities.Atraccion) {
	filePath := "core/database/saves/atracciones_saves.json"

	var prev []entities.Atraccion
	fileContent, _ := os.ReadFile(filePath)
	if len(fileContent) > 0 {
		_ = json.Unmarshal(fileContent, &prev)
	}

	prev = append(prev, data)

	content, err := json.MarshalIndent(prev, "", "  ")
	if err != nil {
		log.Println("❌ Error al serializar atracción para respaldo:", err)
		return
	}

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Println("❌ Error al escribir respaldo de atracción:", err)
	}
}
