package entities

type Atraccion struct {
	Id      int    `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Nombre  string `json:"nombre"`
	Tiempo  int    `json:"tiempo"`
	Hora    string `json:"hora"`
	Fecha   string `json:"fecha"`
	Enviado bool   `json:"enviado"`
	Zona    string `json:"zona"`
}