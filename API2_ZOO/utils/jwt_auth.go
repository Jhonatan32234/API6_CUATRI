package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"time"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("clave_secreta_super_segura")

type CustomClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Zona   string `json:"zona"`

	jwt.RegisteredClaims
}


func ValidateTokenFromQuery(c *gin.Context, allowedRoles ...string) (*CustomClaims, error) {
	tokenString := c.Query("token")
	if tokenString == "" {
		return nil, errors.New("token no proporcionado en la query")
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inválido")
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Println("❌ Error al parsear token:", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Unix() < time.Now().Unix() {
			return nil, errors.New("token expirado")
		}

		if len(allowedRoles) > 0 {
			allowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					allowed = true
					break
				}
			}
			if !allowed {
				return nil, errors.New("acceso denegado: rol no autorizado")
			}
		}

		return claims, nil
	}

	return nil, errors.New("token inválido")
}





