package handlers

import (
	"encoding/json"
	"net/http"
	"portfolio-backend/internal/models"
)

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	projects := []models.Project{
		{ID: "1", Name: "Projeto Rádio Browser", Description: "Uma plataforma de rádio online", TechStack: []string{"React", "Vite", "Tailwind", "Node.js"}},
		{ID: "2", Name: "Portfólio", Description: "Meu site pessoal", TechStack: []string{"Next.js", "Golang", "TypeScript"}},
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(projects)
}
