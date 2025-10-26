package user

type Service struct {
	repo *Repository
}

func NewUserService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(dto DTO) error {
	return s.repo.CreateUser(dto)
}
