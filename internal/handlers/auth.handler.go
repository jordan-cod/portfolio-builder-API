package handlers

import (
	"net/http"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao validar os dados"})
		return
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email já registrado"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar o usuário"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuário cadastrado com sucesso",
		"api_key": user.APIKey,
	})
}
