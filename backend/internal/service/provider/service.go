package provider

import (
	"context"

	"github.com/aivideonote/backend/internal/model"
	"github.com/aivideonote/backend/internal/pkg/logger"
	"github.com/aivideonote/backend/internal/repository"
	"github.com/aivideonote/backend/internal/service/llm"
)

type Service struct {
	repo      *repository.ProviderRepository
	llmClient *llm.Client
}

func NewService(repo *repository.ProviderRepository, llmClient *llm.Client) *Service {
	return &Service{repo: repo, llmClient: llmClient}
}

func (s *Service) Create(provider *model.LLMProvider) error {
	if err := s.repo.Create(provider); err != nil {
		logger.L.Errorf("创建提供商失败: %v", err)
		return err
	}
	return nil
}

func (s *Service) GetAll() ([]model.LLMProvider, error) {
	providers, err := s.repo.FindAll()
	if err != nil {
		logger.L.Errorf("获取提供商列表失败: %v", err)
		return nil, err
	}
	if providers == nil {
		providers = []model.LLMProvider{}
	}
	return providers, nil
}

func (s *Service) GetByID(id string) (*model.LLMProvider, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(id string, updates map[string]interface{}) error {
	if err := s.repo.Update(id, updates); err != nil {
		logger.L.Errorf("更新提供商失败: %v", err)
		return err
	}
	return nil
}

func (s *Service) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		logger.L.Errorf("删除提供商失败: %v", err)
		return err
	}
	return nil
}

func (s *Service) ListModels(ctx context.Context, id string) ([]string, error) {
	provider, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	models, err := s.llmClient.ListModels(ctx, provider.BaseURL, provider.ApiKey)
	if err != nil {
		logger.L.Errorf("获取模型列表失败: %v", err)
		return nil, err
	}
	return models, nil
}

func (s *Service) TestConnection(ctx context.Context, id string) error {
	provider, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	_, err = s.llmClient.ListModels(ctx, provider.BaseURL, provider.ApiKey)
	if err != nil {
		logger.L.Errorf("连接测试失败: %v", err)
		return err
	}
	return nil
}
