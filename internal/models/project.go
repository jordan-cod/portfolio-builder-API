package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Project struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	TechStack   pq.StringArray `json:"techStack" gorm:"type:text[];not null"`
	GitHubUrl   string         `json:"gitHubUrl"`
	LiveUrl     string         `json:"liveUrl"`
	Image       string         `json:"image"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"autoCreateTime; not null; default:now()"`
}
