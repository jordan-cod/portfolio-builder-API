package handlers

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/models"
	repository "portfolio-backend/internal/repositories"
	"strings"

	"github.com/gin-gonic/gin"
)

var projectRepo *repository.ProjectRepository

func getProjectRepo() *repository.ProjectRepository {
	if projectRepo == nil {
		projectRepo = repository.NewProjectRepository(db.DB)
	}
	return projectRepo
}

// Alterar toda lógica de negócios para arquivos de service

// GetAllProjectsHandler godoc
// @Summary      Lista todos os projetos do usuário
// @Description  Retorna todos os projetos com paginação, ordenação e total de registros
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      json
// @Param        page   query     int     false  "Página atual"
// @Param        size   query     int     false  "Tamanho da página"
// @Param        sort   query     string  false  "Campo para ordenar (ex: name, created_at)"
// @Param        order  query     string  false  "Ordem (asc ou desc)"
// @Success 200 {object} models.ProjectListResponse
// @Failure      500    {object}  map[string]string
// @Router       /api/projects [get]
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

// GetOneProjectHandler godoc
// @Summary      Busca um projeto por ID
// @Description  Retorna um projeto específico do usuário autenticado
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id   path      string  true  "ID do projeto"
// @Success      200  {object}  models.ProjectSwagger
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/projects/{id} [get]
func GetOneProjectHandler(c *gin.Context) {
	projectRepo := getProjectRepo()

	projectID := c.Param("id")
	user := c.MustGet("user").(models.User)

	project, err := projectRepo.FindByUserID(projectID, user.ID)

	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Projeto não encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar projeto"})
		}
		log.Println("Erro ao buscar projeto:", err)
		return
	}

	c.JSON(200, project)
}

// TODO: adicionar validações de campos e melhorar retorno de erros

// CreateProjectHandler godoc
// @Summary      Cria um novo projeto
// @Description  Cria um novo projeto para o usuário autenticado
// @Tags         Projects
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project  body      models.ProjectSwagger  true  "Dados do projeto"
// @Success      201      {object}  models.ProjectSwagger
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/projects [post]
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

// UpdateProjectHandler godoc
// @Summary      Atualiza um projeto
// @Description  Atualiza um projeto existente do usuário autenticado
// @Tags         Projects
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "ID do projeto"
// @Param        project  body      models.ProjectSwagger true  "Dados do projeto atualizado"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/projects/{id} [put]
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

// DeleteProjectHandler godoc
// @Summary      Deleta um projeto
// @Description  Remove um projeto do usuário autenticado
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id   path      string  true  "ID do projeto"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/projects/{id} [delete]
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

// FavoriteProjectHandler godoc
// @Summary      Favorita/Desfavorita um projeto
// @Description  Alterna o status de favorito de um projeto do usuário
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id   path      string  true  "ID do projeto"
// @Success      204  {string}  string  "No Content"
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/projects/{id}/favorite [patch]
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

// ExportProjectsToCSVHandler godoc
// @Summary      Exporta projetos para CSV
// @Description  Exporta todos os projetos do usuário autenticado em formato CSV
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      text/csv
// @Success      200  {string}  string  "Arquivo CSV"
// @Failure      500  {object}  map[string]string
// @Router       /api/projects/export/csv [get]
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
