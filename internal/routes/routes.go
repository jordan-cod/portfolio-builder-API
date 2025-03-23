package routes

import (
	"portfolio-backend/internal/handlers"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes(r *chi.Mux) {

	r.Get("/projects", handlers.ProjectsHandler)

}
