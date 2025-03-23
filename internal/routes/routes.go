package routes

import (
	"portfolio-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/projects", handlers.ProjectsHandler)
		api.GET("/projects/:id", handlers.GetProjectHandler)
	}
}
