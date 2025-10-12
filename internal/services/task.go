package services

import (
	"vado_server/internal/models"
	"vado_server/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.repo.GetAll()
}

func (s *TaskService) GetAllByUser(userID uint) ([]models.Task, error) {
	return s.repo.GetAllByUser(userID)
}
