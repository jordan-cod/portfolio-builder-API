package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Check API health
// @Description Returns a simple message to confirm the server is running
// @Tags health
// @Produce json
// @Success      200  {object}  models.HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
