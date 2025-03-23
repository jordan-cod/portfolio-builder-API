package routes

import (
	"portfolio-backend/internal/handlers"
	"portfolio-backend/internal/midlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api", midlewares.ApiKeyMiddleware)
	{
		api.GET("/projects", handlers.ProjectsHandler)
		api.GET("/projects/:id", handlers.GetProjectHandler)
	}
}
