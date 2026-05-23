package repository

import (
	"github.com/google/uuid"
	"github.com/aivideonote/backend/internal/database"
	"github.com/aivideonote/backend/internal/model"
)

type CookieRepository struct{}

func NewCookieRepository() *CookieRepository {
	return &CookieRepository{}
}

func (r *CookieRepository) FindAll() ([]model.PlatformCookie, error) {
	var cookies []model.PlatformCookie
	err := database.DB.Order("platform ASC").Find(&cookies).Error
	if err != nil {
		return nil, err
	}
	if cookies == nil {
		cookies = []model.PlatformCookie{}
	}
	return cookies, nil
}

func (r *CookieRepository) FindByPlatform(platform string) (*model.PlatformCookie, error) {
	var cookie model.PlatformCookie
	err := database.DB.Where("platform = ?", platform).First(&cookie).Error
	if err != nil {
		return nil, err
	}
	return &cookie, nil
}

func (r *CookieRepository) Upsert(platform, content string) (*model.PlatformCookie, error) {
	existing, err := r.FindByPlatform(platform)
	if err == nil {
		existing.Content = content
		if err := database.DB.Save(existing).Error; err != nil {
			return nil, err
		}
		return existing, nil
	}

	cookie := &model.PlatformCookie{
		ID:       uuid.New().String(),
		Platform: platform,
		Content:  content,
	}
	if err := database.DB.Create(cookie).Error; err != nil {
		return nil, err
	}
	return cookie, nil
}

func (r *CookieRepository) Delete(platform string) error {
	return database.DB.Where("platform = ?", platform).Delete(&model.PlatformCookie{}).Error
}
