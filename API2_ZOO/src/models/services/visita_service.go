package services

import (
	"api2/db"
	"api2/src/entities"
	"log"
)

func SaveVisitas(input []entities.Visitas) ([]entities.Visitas, error) {
	var guardadas []entities.Visitas

	for _, item := range input {
		item.Enviado = false
		if err := db.DB.Create(&item).Error; err != nil {
			log.Println("❌ Error al guardar visita:", err)
		} else {
			guardadas = append(guardadas, item)

		}
	}

	if len(guardadas) == 0 {
		log.Println("⚠️ Ninguna visita fue guardada.")
		return nil, nil
	}

	return guardadas, nil
}



func GetVisitaByID(id uint) (*entities.Visitas, error) {
	var visita entities.Visitas
	if err := db.DB.First(&visita, id).Error; err != nil {
		log.Println("❌ Error al obtener visita por ID:", err)
		return nil, err
	}
	return &visita, nil
}
