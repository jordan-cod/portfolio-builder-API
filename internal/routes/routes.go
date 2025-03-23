package routes

import (
	"portfolio-backend/internal/handlers"
	"portfolio-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api", middlewares.ApiKeyMiddleware)
	{
		api.GET("/projects", middlewares.PaginateMiddleware, handlers.GetAllProjectsHandler)
		api.GET("/projects/:id", handlers.GetOneProjectHandler)
	}
}
