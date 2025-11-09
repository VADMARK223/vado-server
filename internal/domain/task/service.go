package task

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllTasks() ([]Task, error) {
	return s.repo.GetAll()
}

func (s *Service) GetAllByUser(userID uint) ([]Task, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *Service) UpdateCompleted(taskID, userID uint, completed bool) error {
	return s.repo.UpdateCompleted(taskID, userID, completed)
}
