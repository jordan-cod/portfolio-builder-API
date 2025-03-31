package repository

import (
	"portfolio-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	*GenericRepository[models.Project]
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{
		GenericRepository: NewGenericRepository[models.Project](db),
	}
}

func (r *ProjectRepository) FindByUserID(projectID string, userID uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.First(&project, "id = ? AND user_id = ?", projectID, userID).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}
