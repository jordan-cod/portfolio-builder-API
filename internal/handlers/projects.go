package handlers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// Alterar toda lógica de negócios para arquivos de service

func GetAllProjectsHandler(c *gin.Context) {
	page := c.MustGet("page").(int)
	limit := c.MustGet("limit").(int)
	sort := c.MustGet("sort").(string)
	order := c.MustGet("order").(string)

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

	if err := query.Order(sort + " " + order).Offset(offset).Limit(limit).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar projetos"})
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

// TODO: adicionar validações de campos e melhorar retorno de erros
func CreateProjectHandler(c *gin.Context) {
	var project models.Project
	user := c.MustGet("user").(models.User)

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		log.Println("Erro ao bindar dados:", err)
		return
	}

	project.UserID = user.ID

	if err := db.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar projeto"})
		log.Println("Erro ao criar projeto:", err)
		return
	}

	c.JSON(http.StatusCreated, project)
}

// TODO: adicionar validações de campos e melhorar retorno de erros
func UpdateProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	var project models.Project
	user := c.MustGet("user").(models.User)

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		log.Println("Erro ao bindar dados:", err)
		return
	}

	result := db.DB.Model(&models.Project{}).Where("id = ? AND user_id = ?", projectID, user.ID).Updates(project)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar projeto"})
		log.Println("Erro ao atualizar projeto:", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Projeto não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Projeto atualizado com sucesso"})
}

func DeleteProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	user := c.MustGet("user").(models.User)

	result := db.DB.Delete(&models.Project{}, "id = ? AND user_id = ?", projectID, user.ID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir projeto"})
		log.Println("Erro ao excluir projeto:", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Projeto não encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Projeto excluído com sucesso"})
}

func FavoriteProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	user := c.MustGet("user").(models.User)

	var project models.Project

	result := db.DB.First(&project, "id = ? AND user_id = ?", projectID, user.ID)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Projeto não encontrado ou não pertence a você"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar projeto"})
		}
		log.Println("Erro ao buscar projeto:", result.Error)
		return
	}

	project.IsFavorited = !project.IsFavorited

	if err := db.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar favorito"})
		log.Println("Erro ao atualizar favorito:", err)
		return
	}

	c.Status(http.StatusNoContent)
}

func ExportProjectsToCSVHandler(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var projects []models.Project
	if err := db.DB.Where("user_id = ?", user.ID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar projetos"})
		log.Println("Erro ao buscar projetos:", err)
		return
	}

	var csvContent []string
	csvContent = append(csvContent, "ID,Nome,Descrição,Status,Repo URL,Tecnologias,Data de Criação")

	for _, project := range projects {
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s",
			project.ID.String(),
			project.Name,
			project.Description,
			project.Status,
			project.GitHubUrl,
			strings.Join(project.TechStack, "|"),
			project.CreatedAt.Format("2006-01-02 15:04:05"),
		)
		csvContent = append(csvContent, line)
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=projetos.csv")
	c.Writer.WriteHeader(http.StatusOK)
	csvWriter := csv.NewWriter(c.Writer)

	for _, line := range csvContent {
		record := strings.Split(line, ",")
		if err := csvWriter.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar CSV"})
			log.Println("Erro ao escrever no CSV:", err)
			return
		}
	}
	csvWriter.Flush()
}
