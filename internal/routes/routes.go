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
			projects.GET("/", middlewares.PaginateMiddleware, handlers.GetAllProjectsHandler)
			projects.GET("/:id", handlers.GetOneProjectHandler)
		}
	}
}
