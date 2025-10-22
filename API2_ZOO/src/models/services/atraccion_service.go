package services

import (
	"api2/db"
	"api2/src/entities"
	"log"
)

func SaveAtracciones(input []entities.Atraccion) ([]entities.Atraccion, error) {
	var guardadas []entities.Atraccion

	for _, item := range input {
		item.Enviado = false
		if err := db.DB.Create(&item).Error; err != nil {
			log.Println("❌ Error al guardar atracción:", err)
		} else {
			guardadas = append(guardadas, item)

		}
	}

	if len(guardadas) == 0 {
		log.Println("⚠️ Ninguna atracción fue guardada.")
		return nil, nil
	}

	return guardadas, nil
}



func GetAtraccionByID(id uint) (*entities.Atraccion, error) {
	var atr entities.Atraccion
	if err := db.DB.First(&atr, id).Error; err != nil {
		log.Println("❌ Error al obtener atracción por ID:", err)
		return nil, err
	}
	return &atr, nil
}
