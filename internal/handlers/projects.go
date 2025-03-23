package handlers

import (
	"log"
	"net/http"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func ProjectsHandler(c *gin.Context) {
	var projects []models.Project

	result := db.DB.Find(&projects)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar projetos"})
		log.Println("Erro ao buscar projetos:", result.Error)
		return
	}

	c.JSON(http.StatusOK, projects)
}

func GetProjectHandler(c *gin.Context) {
	projectID := c.Param("id")

	var project models.Project

	result := db.DB.First(&project, "id = ?", projectID)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Projeto não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar projeto"})
		}
		log.Println("Erro ao buscar projeto:", result.Error)
		return
	}

	c.JSON(200, project)
}
