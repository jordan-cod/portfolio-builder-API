package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Status string

const (
	StatusConcluido   Status = "Concluído"
	StatusEmAndamento Status = "Em andamento"
	StatusNaoIniciado Status = "Não iniciado"
)

type Project struct {
	ID          uuid.UUID      `json:"id" gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	TechStack   pq.StringArray `json:"techStack" gorm:"type:text[];not null"`
	Status      Status         `json:"status" gorm:"not null;default:'Concluído'"`
	GitHubUrl   string         `json:"gitHubUrl"`
	LiveUrl     string         `json:"liveUrl"`
	Image       string         `json:"image"`
	IsFavorited bool           `json:"isFavorited" gorm:"default:false;not null"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"autoCreateTime; not null; default:now()"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"autoUpdateTime; not null; default:now()"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User   User      `json:"-" gorm:"foreignKey:UserID;references:ID"`
}

func (project *Project) BeforeCreate(tx *gorm.DB) (err error) {
	if project.Status != StatusConcluido && project.Status != StatusEmAndamento && project.Status != StatusNaoIniciado {
		return fmt.Errorf("status inválido: %s. O status deve ser 'Concluído', 'Em andamento' ou 'Não iniciado'", project.Status)
	}
	return nil
}

func (project *Project) BeforeUpdate(tx *gorm.DB) (err error) {
	if project.Status != StatusConcluido && project.Status != StatusEmAndamento && project.Status != StatusNaoIniciado {
		return fmt.Errorf("status inválido: %s. O status deve ser 'Concluído', 'Em andamento' ou 'Não iniciado'", project.Status)
	}
	return nil
}
