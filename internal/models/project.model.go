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

type ProjectSwagger struct {
	ID          string   `json:"id" example:"1"`
	Title       string   `json:"title" example:"Portfolio Builder"`
	Description string   `json:"description" example:"Um criador de portfólios para devs"`
	TechStack   []string `json:"techStack" example:"React,TypeScript,Node.js"`
	Link        string   `json:"link" example:"https://meuportfolio.com"`
	Image       string   `json:"image" example:"https://meuportfolio.com/capa.png"`
	IsFavorite  bool     `json:"isFavorite" example:"false"`
	CreatedAt   string   `json:"createdAt" example:"2025-04-04T12:00:00Z"`
	UpdatedAt   string   `json:"updatedAt" example:"2025-04-04T12:00:00Z"`
}

type ProjectListResponse struct {
	Data       []ProjectSwagger `json:"data"`
	TotalCount int              `json:"totalCount"`
	Page       int              `json:"page"`
	Size       int              `json:"size"`
}

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
