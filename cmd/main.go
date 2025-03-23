package main

import (
	"log"
	"net/http"

	"portfolio-backend/internal/routes"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	routes.SetupRoutes(r)

	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
