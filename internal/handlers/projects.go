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
// @Summary      Get all projects
// @Description  Returns all user's projects with pagination and sorting
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      json
// @Param        page   query     int     false  "Current page"
// @Param        size   query     int     false  "Page size"
// @Param        sort   query     string  false  "Sort field (e.g. name, created_at)"
// @Param        order  query     string  false  "Sort order (asc or desc)"
// @Success      200    {object}  models.ProjectListResponse
// @Failure      500    {object}  map[string]string
// @Router       /projects [get]
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count projects"})
		log.Println("Count error:", err)
		return
	}

	offset := (page - 1) * limit

	if err := query.Order(sort + " " + order).Offset(offset).Limit(limit).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
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
// @Summary      Get project by ID
// @Description  Returns a specific project by ID for the authenticated user
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  models.ProjectSwagger
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id} [get]
func GetOneProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	user := c.MustGet("user").(models.User)

	project, err := getProjectRepo().FindByUserID(projectID, user.ID)

	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		}
		log.Println("Find error:", err)
		return
	}

	c.JSON(http.StatusOK, project)
}

// TODO: adicionar validações de campos e melhorar retorno de erros

// CreateProjectHandler godoc
// @Summary      Create a project
// @Description  Creates a new project for the authenticated user
// @Tags         Projects
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        project  body      models.ProjectSwagger  true  "Project data"
// @Success      201      {object}  models.ProjectSwagger
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /projects [post]
func CreateProjectHandler(c *gin.Context) {
	var project models.Project
	user := c.MustGet("user").(models.User)

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		log.Println("Bind error:", err)
		return
	}

	project.UserID = user.ID

	if err := db.DB.Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		log.Println("Create error:", err)
		return
	}

	c.JSON(http.StatusCreated, project)
}

// TODO: adicionar validações de campos e melhorar retorno de erros

// UpdateProjectHandler godoc
// @Summary      Update a project
// @Description  Updates a project of the authenticated user
// @Tags         Projects
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Project ID"
// @Param        project  body      models.ProjectSwagger true  "Updated project data"
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /projects/{id} [put]
func UpdateProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	var project models.Project
	user := c.MustGet("user").(models.User)

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		log.Println("Bind error:", err)
		return
	}

	result := db.DB.Model(&models.Project{}).Where("id = ? AND user_id = ?", projectID, user.ID).Updates(project)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		log.Println("Update error:", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

// DeleteProjectHandler godoc
// @Summary      Delete a project
// @Description  Deletes a project from the authenticated user
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id} [delete]
func DeleteProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	user := c.MustGet("user").(models.User)

	result := db.DB.Delete(&models.Project{}, "id = ? AND user_id = ?", projectID, user.ID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		log.Println("Delete error:", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// FavoriteProjectHandler godoc
// @Summary      Toggle favorite
// @Description  Toggles the favorite status of a project
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      json
// @Param        id   path      string  true  "Project ID"
// @Success      204  {string}  string  "No Content"
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /projects/{id}/favorite [patch]
func FavoriteProjectHandler(c *gin.Context) {
	projectID := c.Param("id")
	user := c.MustGet("user").(models.User)

	var project models.Project

	result := db.DB.First(&project, "id = ? AND user_id = ?", projectID, user.ID)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found or does not belong to you"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project"})
		}
		log.Println("Find error:", result.Error)
		return
	}

	project.IsFavorited = !project.IsFavorited

	if err := db.DB.Save(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle favorite"})
		log.Println("Favorite error:", err)
		return
	}

	c.Status(http.StatusNoContent)
}

// ExportProjectsToCSVHandler godoc
// @Summary      Export projects to CSV
// @Description  Exports all authenticated user's projects in CSV format
// @Tags         Projects
// @Security     ApiKeyAuth
// @Produce      text/csv
// @Success      200  {string}  string  "CSV file"
// @Failure      500  {object}  map[string]string
// @Router       /projects/export/csv [get]
func ExportProjectsToCSVHandler(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var projects []models.Project
	if err := db.DB.Where("user_id = ?", user.ID).Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		log.Println("Fetch error:", err)
		return
	}

	var csvContent []string
	csvContent = append(csvContent, "ID,Name,Description,Status,Repo URL,Tech Stack,Created At")

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
	c.Header("Content-Disposition", "attachment;filename=projects.csv")
	c.Writer.WriteHeader(http.StatusOK)
	csvWriter := csv.NewWriter(c.Writer)

	for _, line := range csvContent {
		record := strings.Split(line, ",")
		if err := csvWriter.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV"})
			log.Println("CSV write error:", err)
			return
		}
	}
	csvWriter.Flush()
}
