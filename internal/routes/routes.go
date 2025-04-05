package routes

import (
	"portfolio-backend/internal/handlers"
	"portfolio-backend/internal/middlewares"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(middlewares.CORSConfig())

	api := r.Group("/api")
	{

		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		api.POST("/register", handlers.Register)
		api.GET("/health", handlers.HealthCheck)

		protected := api.Group("/", middlewares.APIKeyAuthMiddleware())
		{
			protected.POST("/user/api-key", handlers.RenewAPIKey)

			projects := protected.Group("/projects")
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
}
