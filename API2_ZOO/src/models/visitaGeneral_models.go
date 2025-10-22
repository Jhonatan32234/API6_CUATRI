package models

import (
	"api2/db"
	"api2/src/entities"
	"errors"
)

func CreateVisitaGeneral(visita entities.VisitaGeneral) error {
	var existing entities.VisitaGeneral
	result := db.DB.Where("fecha = ?", visita.Fecha).First(&existing)
	if result.Error == nil {
		return errors.New("ya existe una visita con esa fecha")
	}

	return db.DB.Create(&visita).Error
}


func GetAllVisitasGeneral() ([]entities.VisitaGeneral, error) {
	var visitas []entities.VisitaGeneral
	err := db.DB.Find(&visitas).Error
	return visitas, err
}

func GetVisitaGeneralByFecha(fecha string, visita *entities.VisitaGeneral) error {
	return db.DB.Where("fecha = ?", fecha).First(visita).Error
}


func UpdateVisitaGeneralPorFecha(fecha string, updated entities.VisitaGeneral) error {
	var visita entities.VisitaGeneral

	// Buscar la visita por fecha
	err := db.DB.Where("fecha = ?", fecha).First(&visita).Error
	if err != nil {
		return err
	}

	// Validar que no exista otra visita con la misma fecha si actualizas la fecha
	if updated.Fecha != "" && updated.Fecha != fecha {
		var existing entities.VisitaGeneral
		result := db.DB.Where("fecha = ? AND id != ?", updated.Fecha, visita.Id).First(&existing)
		if result.Error == nil {
			return errors.New("ya existe otra visita con esa fecha")
		}
		visita.Fecha = updated.Fecha
	}

	// Actualizar otros campos (ajusta según tus campos)
	if updated.Visitas != 0 {
		visita.Visitas = updated.Visitas
	}

	return db.DB.Save(&visita).Error
}


func DeleteVisitaGeneralPorFecha(fecha string) error {
	return db.DB.Where("fecha = ?", fecha).Delete(&entities.VisitaGeneral{}).Error
}

