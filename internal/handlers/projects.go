package handlers

import (
	"log"
	"net/http"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func GetAllProjectsHandler(c *gin.Context) {
	page := c.MustGet("page").(int)
	limit := c.MustGet("limit").(int)

	user := c.MustGet("user").(models.User)

	var projects []models.Project
	var totalCount int64

	query := db.DB.Model(&models.Project{}).Where("user_id = ?", user.ID)

	if err := query.Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar projetos"})
		log.Println("Erro ao contar projetos:", err)
		return
	}

	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar projetos"})
		log.Println("Erro ao buscar projetos:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       projects,
		"totalCount": totalCount,
		"page":       page,
		"size":       limit,
	})
}

func GetOneProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	user := c.MustGet("user").(models.User)

	var project models.Project

	result := db.DB.First(&project, "id = ? AND user_id = ?", projectID, user.ID)
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
