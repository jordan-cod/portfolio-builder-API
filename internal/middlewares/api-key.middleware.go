package middlewares

import (
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Chave de API não fornecida"})
			return
		}

		var user models.User
		if err := db.DB.Where("api_key = ?", apiKey).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Chave de API inválida"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
