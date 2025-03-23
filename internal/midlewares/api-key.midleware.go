package midlewares

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware(c *gin.Context) {
	apiKey := c.GetHeader("X-API-Key")

	if apiKey != os.Getenv("API_KEY") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}
