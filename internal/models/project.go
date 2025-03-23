package models

import (
	_ "gorm.io/gorm"
)

type Project struct {
	ID          string   `json:"id" gorm:"primaryKey"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	TechStack   []string `json:"techStack" gorm:"type:text[]"`
	GitHubUrl   string   `json:"gitHubUrl"`
	LiveUrl     string   `json:"liveUrl"`
	Image       string   `json:"image"`
	CreatedAt   string   `json:"createdAt"`
}

func (Project) TableName() string {
	return "projects"
}
