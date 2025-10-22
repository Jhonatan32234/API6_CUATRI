package models

import (
	"encoding/base64"
	"net/http"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`           // omitido al retornar
	Role     string `json:"role"`
	Zona     string `json:"zona"`
	Image    []byte `json:"-"`           // imagen en crudo (no se expone en JSON)
	ImageStr string `json:"image"`       // imagen codificada base64
	MimeType string `json:"imageType"`   // tipo MIME (ej: image/png)
}

// Procesa la imagen para mostrarla en JSON
func (u *User) FormatImage() {
	if len(u.Image) > 0 {
		u.MimeType = http.DetectContentType(u.Image)
		u.ImageStr = base64.StdEncoding.EncodeToString(u.Image)
	}
}
