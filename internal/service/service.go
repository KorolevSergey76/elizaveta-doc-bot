package service

import (
	"cosmetologybotliza/internal/domain"
	"cosmetologybotliza/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveUser(id int64, username, firstName string) error {
	return s.repo.SaveUser(id, username, firstName)
}

func (s *Service) GetServices() ([]domain.Service, error) {
	return s.repo.GetServices()
}

func (s *Service) GetMasters() ([]domain.Master, error) {
	return s.repo.GetMasters()
}
