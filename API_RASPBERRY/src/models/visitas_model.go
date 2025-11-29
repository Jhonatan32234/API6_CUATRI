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

func GetVisitasFromDate(fecha string) ([]entities.Visitas, error) {
	var visitas []entities.Visitas

	// Obtener visitas desde la fecha que NO hayan sido enviadas (opcional)
	err := database.DB.Where("fecha >= ?", fecha).Find(&visitas).Error
	if err != nil {
		return nil, err
	}

	if len(visitas) == 0 {
		log.Println("No hay visitas para enviar desde la fecha:", fecha)
		return visitas, nil
	}

	// Publicar en RabbitMQ
	if rabbitmq.PublishToTopic(visitas, "visitas_topic", "visita.data") {
		// Marcar como enviadas
		err = database.DB.Model(&entities.Visitas{}).
			Where("fecha >= ? AND enviado = ?", fecha, false).
			Update("enviado", true).Error
		if err != nil {
			log.Println("Error actualizando visitas como enviadas:", err)
		}
	}

	return visitas, nil
}

func validateVisita(visita entities.Visitas) error {
	// Validar que Visitantes sea mayor a 0
	if visita.Visitantes <= 0 {
		return fmt.Errorf("el campo 'visitantes' debe ser mayor a 0")
	}

	// Validar que no haya campos vacíos o nulos
	if visita.Hora == "" || strings.TrimSpace(visita.Hora) == "" {
		return fmt.Errorf("el campo 'hora' es requerido")
	}

	if visita.Fecha == "" || strings.TrimSpace(visita.Fecha) == "" {
		return fmt.Errorf("el campo 'fecha' es requerido")
	}

	if visita.Zona == "" || strings.TrimSpace(visita.Zona) == "" {
		return fmt.Errorf("el campo 'zona' es requerido")
	}

	return nil
}

func SaveVisitas(input []entities.Visitas) ([]entities.Visitas, error) {
	var guardadas []entities.Visitas
	var errores []string

	for _, item := range input {
		// Validar cada visita antes de guardar
		if err := validateVisita(item); err != nil {
			errores = append(errores, fmt.Sprintf("Visita inválida: %v", err))
			log.Printf("❌ Validación fallida para visita: %v", err)
			continue
		}

		item.Enviado = false
		if err := database.DB.Create(&item).Error; err != nil {
			log.Println("❌ Error al guardar visita:", err)
			saveVisitaToFile(item)
		} else {
			guardadas = append(guardadas, item)
		}
	}

	// Si hay errores de validación, retornarlos
	if len(errores) > 0 && len(guardadas) == 0 {
		return nil, fmt.Errorf("errores de validación: %s", strings.Join(errores, "; "))
	}

	if len(guardadas) == 0 {
		log.Println("⚠️ Ninguna visita fue guardada. No se enviará al broker.")
		return nil, nil
	}

	var toSend []entities.Visitas
	database.DB.Where("enviado = ?", false).Find(&toSend)

	if len(toSend) > 0 && rabbitmq.PublishToTopic(toSend, "visitas_topic", "visita.data") {
		database.DB.Model(&entities.Visitas{}).Where("enviado = ?", false).Update("enviado", true)

		// 🔁 Enviar por zona solo el ID
		for _, item := range toSend {
			rabbitmq.PublishIDToZoneTopic("visitas_topic", item.Zona, item.Id, "visitas")
		}
	}

	return toSend, nil
}


func saveVisitaToFile(data entities.Visitas) {
	filePath := "core/database/saves/visitas_saves.json"

	var prev []entities.Visitas
	fileContent, _ := os.ReadFile(filePath)
	if len(fileContent) > 0 {
		_ = json.Unmarshal(fileContent, &prev)
	}

	prev = append(prev, data)

	content, err := json.MarshalIndent(prev, "", "  ")
	if err != nil {
		log.Println("❌ Error al serializar visita para respaldo:", err)
		return
	}

	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		log.Println("❌ Error al escribir respaldo de visita:", err)
	}
}
