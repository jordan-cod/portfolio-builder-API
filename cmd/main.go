package main

import (
	"log"
	"net/http"

	"portfolio-backend/internal/config"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/routes"

	"github.com/go-chi/chi/v5"
)

func main() {
	config.LoadEnv()

	err := db.ConnectBD()
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}

	r := chi.NewRouter()

	routes.SetupRoutes(r)

	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
