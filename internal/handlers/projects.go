package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/models"
)

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	var projects []models.Project

	result := db.DB.Find(&projects)
	if result.Error != nil {
		http.Error(w, "Erro ao buscar projetos", http.StatusInternalServerError)
		log.Println("Erro ao buscar projetos:", result.Error)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(projects)
	if err != nil {
		http.Error(w, "Erro ao codificar resposta", http.StatusInternalServerError)
		log.Println("Erro ao codificar resposta:", err)
	}
}
