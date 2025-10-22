package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTQueryMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := ValidateTokenFromQuery(c, roles...)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("zona", claims.Zona) 

		c.Next()
	}
}


