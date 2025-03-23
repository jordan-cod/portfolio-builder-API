package routes

import (
	"portfolio-backend/internal/handlers"
	"portfolio-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	r.Use(middlewares.CORSConfig())

	api := r.Group("/api", middlewares.ApiKeyMiddleware)
	{
		api.GET("/projects", middlewares.PaginateMiddleware, handlers.GetAllProjectsHandler)
		api.GET("/projects/:id", handlers.GetOneProjectHandler)

		api.POST("/auth/register", handlers.Register)
	}
}
