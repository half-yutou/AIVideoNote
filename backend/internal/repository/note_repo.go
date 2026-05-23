package repository

import (
	"github.com/google/uuid"
	"github.com/aivideonote/backend/internal/database"
	"github.com/aivideonote/backend/internal/model"
)

type NoteRepository struct{}

func NewNoteRepository() *NoteRepository {
	return &NoteRepository{}
}

func (r *NoteRepository) Create(note *model.NoteRecord) error {
	if note.ID == "" {
		note.ID = uuid.New().String()
	}
	return database.DB.Create(note).Error
}

func (r *NoteRepository) FindByTaskID(taskID string) (*model.NoteRecord, error) {
	var note model.NoteRecord
	err := database.DB.Where("task_id = ?", taskID).Order("version DESC").First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepository) FindByID(id string) (*model.NoteRecord, error) {
	var note model.NoteRecord
	err := database.DB.Where("id = ?", id).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *NoteRepository) DeleteByTaskID(taskID string) error {
	return database.DB.Delete(&model.NoteRecord{}, "task_id = ?", taskID).Error
}
