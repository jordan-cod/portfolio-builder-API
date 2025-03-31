package middlewares

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PaginateMiddleware(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("size", "10")
	sort := c.DefaultQuery("sort", "created_at")
	order := c.DefaultQuery("order", "desc")

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

	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"name":       true,
		"status":     true,
	}

	if !allowedSortFields[sort] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campo de ordenação inválido"})
		return
	}

	if order != "asc" && order != "desc" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ordem de ordenação inválida"})
		return
	}

	c.Set("page", page)
	c.Set("limit", limit)
	c.Set("sort", sort)
	c.Set("order", order)

	c.Next()
}
