package gorm

import (
	"vado_server/internal/domain/task"

	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) task.Repository {
	return &TaskRepo{db: db}
}

func (r TaskRepo) GetAll() ([]task.Task, error) {
	var entities []TaskEntity
	err := r.db.Find(&entities).Error

	result := make([]task.Task, 0, len(entities))
	for _, entity := range entities {
		result = append(result, task.Task{
			ID:   entity.ID,
			Name: entity.Name,
		})
	}

	return result, err
}

func (r TaskRepo) GetAllByUserID(ID uint) ([]task.Task, error) {
	var entities []TaskEntity
	err := r.db.Where("user_id = ?", ID).Find(&entities).Error

	result := make([]task.Task, 0, len(entities))
	for _, entity := range entities {
		result = append(result, task.Task{
			ID:          entity.ID,
			Name:        entity.Name,
			Description: entity.Description,
			Completed:   entity.Completed,
			CreatedAt:   entity.CreatedAt,
		})
	}

	return result, err
}
