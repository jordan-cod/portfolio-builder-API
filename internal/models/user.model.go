package models

import (
	"portfolio-backend/internal/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null; unique"`
	Password  string    `json:"password" gorm:"not null"`
	APIKey    string    `json:"api_key" gorm:"not null; unique"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime; not null; default:now()"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime; not null; default:now()"`
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	user.APIKey, err = utils.GenerateAPIKey()
	if err != nil {
		return err
	}

	return nil
}
