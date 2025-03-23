package handlers

import (
	"log"
	"net/http"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllProjectsHandler(c *gin.Context) {

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

	var projects []models.Project
	var totalCount int64

	if err := db.DB.Model(&models.Project{}).Count(&totalCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao contar projetos"})
		log.Println("Erro ao contar projetos:", err)
		return
	}
	offset := (page - 1) * limit

	result := db.DB.Offset(offset).Limit(limit).Find(&projects)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar projetos"})
		log.Println("Erro ao buscar projetos:", result.Error)
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
