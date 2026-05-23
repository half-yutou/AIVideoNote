package cookie

import (
	"github.com/aivideonote/backend/internal/model"
	"github.com/aivideonote/backend/internal/repository"
)

type Service struct {
	repo *repository.CookieRepository
}

func NewService(repo *repository.CookieRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAll() ([]model.PlatformCookie, error) {
	return s.repo.FindAll()
}

func (s *Service) Save(platform, content string) (*model.PlatformCookie, error) {
	return s.repo.Upsert(platform, content)
}

func (s *Service) Delete(platform string) error {
	return s.repo.Delete(platform)
}

func (s *Service) FindByPlatform(platform string) (*model.PlatformCookie, error) {
	return s.repo.FindByPlatform(platform)
}
