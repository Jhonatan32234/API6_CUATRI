package entities

type Visitas struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Visitantes int   `json:"visitantes"`
	Hora      string `json:"hora"`
	Fecha     string `json:"fecha"`
	Zona      string `json:"zona"` // Nuevo campo
	Enviado   bool   `json:"enviado"`
}

type Atraccion struct {
	Id      int    `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Nombre  string `json:"nombre"`
	Tiempo  int    `json:"tiempo"`
	Hora    string `json:"hora"`
	Fecha   string `json:"fecha"`
	Zona    string `json:"zona"` // Nuevo campo
	Enviado bool   `json:"enviado"`
}

