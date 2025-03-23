package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null"`
	Password  string    `json:"password" gorm:"not null"`
	APIKey    string    `json:"api_key" gorm:"not null; unique"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime; not null; default:now()"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime; not null; default:now()"`
}
