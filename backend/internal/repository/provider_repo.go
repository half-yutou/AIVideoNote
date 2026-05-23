package repository

import (
	"github.com/google/uuid"
	"github.com/aivideonote/backend/internal/database"
	"github.com/aivideonote/backend/internal/model"
	"github.com/aivideonote/backend/internal/pkg/crypto"
)

type ProviderRepository struct{}

func NewProviderRepository() *ProviderRepository {
	return &ProviderRepository{}
}

func (r *ProviderRepository) Create(provider *model.LLMProvider) error {
	if provider.ID == "" {
		provider.ID = uuid.New().String()
	}
	r.encrypt(provider)
	return database.DB.Create(provider).Error
}

func (r *ProviderRepository) FindByID(id string) (*model.LLMProvider, error) {
	var provider model.LLMProvider
	err := database.DB.Where("id = ?", id).First(&provider).Error
	if err != nil {
		return nil, err
	}
	r.decrypt(&provider)
	return &provider, nil
}

func (r *ProviderRepository) FindAll() ([]model.LLMProvider, error) {
	var providers []model.LLMProvider
	err := database.DB.Order("created_at DESC").Find(&providers).Error
	if err != nil {
		return nil, err
	}
	for i := range providers {
		r.decrypt(&providers[i])
	}
	return providers, err
}

func (r *ProviderRepository) FindAllEnabled() ([]model.LLMProvider, error) {
	var providers []model.LLMProvider
	err := database.DB.Where("enabled = ?", true).Order("created_at DESC").Find(&providers).Error
	if err != nil {
		return nil, err
	}
	for i := range providers {
		r.decrypt(&providers[i])
	}
	return providers, err
}

func (r *ProviderRepository) Update(id string, updates map[string]interface{}) error {
	if apiKey, ok := updates["api_key"]; ok {
		if s, ok := apiKey.(string); ok && s != "" {
			encrypted, err := crypto.Encrypt(s)
			if err != nil {
				return err
			}
			updates["api_key"] = encrypted
		}
	}
	return database.DB.Model(&model.LLMProvider{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ProviderRepository) Delete(id string) error {
	return database.DB.Delete(&model.LLMProvider{}, "id = ?", id).Error
}

func (r *ProviderRepository) encrypt(p *model.LLMProvider) {
	if p.ApiKey == "" {
		return
	}
	encrypted, err := crypto.Encrypt(p.ApiKey)
	if err != nil {
		return
	}
	p.ApiKey = encrypted
}

func (r *ProviderRepository) decrypt(p *model.LLMProvider) {
	if p.ApiKey == "" {
		return
	}
	decrypted, err := crypto.Decrypt(p.ApiKey)
	if err != nil {
		return
	}
	p.ApiKey = decrypted
}
