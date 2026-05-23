package repository

import (
	"github.com/google/uuid"
	"github.com/aivideonote/backend/internal/database"
	"github.com/aivideonote/backend/internal/model"
)

type GroupRepository struct{}

func NewGroupRepository() *GroupRepository {
	return &GroupRepository{}
}

func (r *GroupRepository) Create(name string) (*model.Group, error) {
	g := &model.Group{
		ID:   uuid.New().String(),
		Name: name,
	}
	if err := database.DB.Create(g).Error; err != nil {
		return nil, err
	}
	return g, nil
}

func (r *GroupRepository) FindAll() ([]model.Group, error) {
	var groups []model.Group
	err := database.DB.Order("created_at ASC").Find(&groups).Error
	return groups, err
}

func (r *GroupRepository) FindByID(id string) (*model.Group, error) {
	var g model.Group
	err := database.DB.Where("id = ?", id).First(&g).Error
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GroupRepository) UpdateName(id, name string) error {
	return database.DB.Model(&model.Group{}).Where("id = ?", id).Update("name", name).Error
}

func (r *GroupRepository) Delete(id string) error {
	return database.DB.Delete(&model.Group{}, "id = ?", id).Error
}

func (r *GroupRepository) EnsureDefault() (*model.Group, error) {
	var g model.Group
	err := database.DB.Where("id = ?", "default").First(&g).Error
	if err == nil {
		return &g, nil
	}
	g = model.Group{
		ID:   "default",
		Name: "默认",
	}
	if err := database.DB.Create(&g).Error; err != nil {
		return nil, err
	}
	return &g, nil
}
