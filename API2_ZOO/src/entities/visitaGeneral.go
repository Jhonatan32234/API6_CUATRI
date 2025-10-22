package entities

type VisitaGeneral struct {
	Id      int    `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Fecha   string `json:"fecha"`
	Visitas int    `json:"visitas"`
}
