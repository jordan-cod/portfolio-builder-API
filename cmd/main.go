package main

import (
	"log"
	"net/http"

	"portfolio-backend/internal/config"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/routes"

	"github.com/gin-gonic/gin"

	_ "portfolio-backend/docs"
)

// @title Portfolio Builder API
// @version 1.0
// @description Uma API open source para o gerenciamento de portfólios de desenvolvedores.
// @host localhost:8080
// @BasePath /api
func main() {
	config.LoadEnv()

	err := db.ConnectBD()
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}

	r := gin.Default()

	routes.SetupRoutes(r)

	log.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
