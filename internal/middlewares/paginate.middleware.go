package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PaginateMiddleware(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)

	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Página inválida"})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tamanho inválido"})
		return
	}

	c.Set("page", page)
	c.Set("limit", limit)

	c.Next()
}
