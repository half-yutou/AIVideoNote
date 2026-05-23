package group

import (
	"github.com/aivideonote/backend/internal/model"
	"github.com/aivideonote/backend/internal/repository"
)

type Service struct {
	repo *repository.GroupRepository
}

func NewService(repo *repository.GroupRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List() ([]model.Group, error) {
	return s.repo.FindAll()
}

func (s *Service) Create(name string) (*model.Group, error) {
	return s.repo.Create(name)
}

func (s *Service) Rename(id, name string) error {
	return s.repo.UpdateName(id, name)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
