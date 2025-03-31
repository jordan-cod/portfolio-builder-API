package routes

import (
	"portfolio-backend/internal/handlers"
	"portfolio-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middlewares.CORSConfig())

	api := r.Group("/api")
	{
		api.POST("/register", handlers.Register)
		api.GET("/health", handlers.HealthCheck)

		projects := api.Group("/projects", middlewares.APIKeyAuthMiddleware())
		{
			projects.POST("/", handlers.CreateProjectHandler)
			projects.GET("/", middlewares.PaginateMiddleware, handlers.GetAllProjectsHandler)
			projects.GET("/:id", handlers.GetOneProjectHandler)
			projects.PUT("/:id", handlers.UpdateProjectHandler)
			projects.DELETE("/:id", handlers.DeleteProjectHandler)

			projects.PATCH("/:id/favorite", handlers.FavoriteProjectHandler)

			projects.GET("/export/csv", handlers.ExportProjectsToCSVHandler)
		}
	}
}
