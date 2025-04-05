package handlers

import (
	"log"
	"net/http"
	"portfolio-backend/internal/db"
	"portfolio-backend/internal/models"
	"portfolio-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Registers a new user
// @Description Creates a new user with name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param registerData body models.RegisterInput true "User registration data"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	var existingUser models.User
	if err := db.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User successfully registered",
		"api_key": user.APIKey,
	})
}

// RenewAPIKey godoc
// @Summary Rotate the user's API key
// @Description Generates a new API key for the authenticated user and replaces the old one
// @Tags auth
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]string "API key successfully rotated"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /user/api-key [post]
func RenewAPIKey(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	newKey, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new API key"})
		log.Println("Error generating new API key:", err)
		return
	}

	result := db.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("api_key", newKey)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update API key"})
		log.Println("Error updating API key:", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key successfully rotated",
		"api_key": newKey,
	})
}
