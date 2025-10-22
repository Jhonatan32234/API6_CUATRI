package entities

type Visitas struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Visitantes int   `json:"visitantes"`
	Hora      string `json:"hora"`
	Fecha     string `json:"fecha"`
	Enviado   bool   `json:"enviado"`
	Zona	  string `json:"zona"`
}