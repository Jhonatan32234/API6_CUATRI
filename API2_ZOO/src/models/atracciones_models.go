package models

import (
	"api2/db"
)


type NowAtraccion struct {
	Fecha  string `json:"fecha"`
	Nombre string `json:"nombre"`
	Hora   string `json:"hora"`
	Total  int    `json:"total"`
}

type LastWeekAtraccion struct {
	Fecha  string `json:"fecha"`
	Nombre string `json:"nombre"`
	Total  int    `json:"total"`
}

type YesterdayAtraccion struct {
	Fecha  string `json:"fecha"`
	Nombre string `json:"nombre"`
	Zona   string `json:"zona"`
	Total  int    `json:"total"`
}

type OjivaAtraccion struct {
	Fecha string `json:"fecha"`
	Hora  string `json:"hora"`
	Total int    `json:"total"`
}

func GetNowAtraccion(zona string) ([]NowAtraccion, error) {
	var result []NowAtraccion
	err := db.DB.Raw(`
		WITH datos AS (
			SELECT 
				fecha,
				nombre,
				CAST(SUBSTRING(hora, 1, 2) AS UNSIGNED) AS hora_truncada,
				SUM(tiempo) AS total
			FROM atraccion
			WHERE fecha = (
				SELECT MAX(fecha) FROM atraccion WHERE zona = ?
			)
			AND zona = ?
			AND CAST(SUBSTRING(hora, 1, 2) AS UNSIGNED) BETWEEN 9 AND 16
			GROUP BY fecha, nombre, CAST(SUBSTRING(hora, 1, 2) AS UNSIGNED)
		),
		acumulado AS (
			SELECT 
				fecha,
				nombre,
				hora_truncada,
				SUM(total) OVER (PARTITION BY nombre ORDER BY hora_truncada) AS total
			FROM datos
		)
		SELECT 
			fecha,
			nombre,
			CONCAT(LPAD(hora_truncada, 2, '0'), ':00') AS hora,
			total
		FROM acumulado
		ORDER BY nombre, hora_truncada
	`, zona, zona).Scan(&result).Error
	return result, err
}


func GetLastWeekAtraccion(zona string) ([]LastWeekAtraccion, error) {
	var fechas []string
	err := db.DB.Raw(`
		SELECT DISTINCT fecha 
		FROM atraccion 
		WHERE zona = ?
		ORDER BY fecha DESC 
		LIMIT 6
	`, zona).Scan(&fechas).Error
	if err != nil || len(fechas) == 0 {
		return nil, err
	}

	var result []LastWeekAtraccion
	err = db.DB.Raw(`
		SELECT fecha, nombre, SUM(tiempo) as total
		FROM atraccion
		WHERE fecha IN (?) AND zona = ? AND CAST(SUBSTRING(hora, 1, 2) AS UNSIGNED) BETWEEN 9 AND 16
		GROUP BY fecha, nombre
		ORDER BY fecha DESC
	`, fechas, zona).Scan(&result).Error
	return result, err
}

func GetYesterdayAtraccion(zona string) ([]YesterdayAtraccion, error) {
	var fecha string
	err := db.DB.Raw(`
		SELECT DISTINCT fecha 
		FROM atraccion 
		WHERE zona = ?
		ORDER BY fecha DESC
		LIMIT 1 OFFSET 1
	`, zona).Scan(&fecha).Error
	if err != nil || fecha == "" {
		return nil, err
	}

	var result []YesterdayAtraccion
	err = db.DB.Raw(`
		SELECT 
			fecha,
			nombre,
			zona,
			SUM(tiempo) AS total
		FROM atraccion
		WHERE 
			fecha = ? 
			AND zona = ? 
			AND CAST(SUBSTRING(hora, 1, 2) AS UNSIGNED) BETWEEN 9 AND 16
		GROUP BY fecha, nombre, zona
	`, fecha, zona).Scan(&result).Error

	return result, err
}

func GetOjivaAtraccion(fecha, zona string) ([]OjivaAtraccion, error) {
	var result []OjivaAtraccion

	if fecha == "" {
		err := db.DB.Raw(`SELECT MAX(fecha) FROM atraccion WHERE zona = ?`, zona).Scan(&fecha).Error
		if err != nil {
			return result, err
		}
	}

	err := db.DB.Raw(`
	SELECT
		fecha,
		DATE_FORMAT(STR_TO_DATE(hora, '%H:%i'), '%H:00') AS hora,
		SUM(tiempo) AS total
	FROM atraccion
	WHERE
		fecha = ? AND
		zona = ? AND
		HOUR(STR_TO_DATE(hora, '%H:%i')) BETWEEN 9 AND 16
	GROUP BY
		fecha, DATE_FORMAT(STR_TO_DATE(hora, '%H:%i'), '%H:00')
	ORDER BY hora
`, fecha, zona).Scan(&result).Error



	return result, err
}
