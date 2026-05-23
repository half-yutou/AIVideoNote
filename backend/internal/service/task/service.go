package task

import (
	"github.com/aivideonote/backend/internal/model"
	"github.com/aivideonote/backend/internal/repository"
)

type Service struct {
	taskRepo  *repository.TaskRepository
	noteRepo  *repository.NoteRepository
	groupRepo *repository.GroupRepository
}

func NewService(taskRepo *repository.TaskRepository, noteRepo *repository.NoteRepository, groupRepo *repository.GroupRepository) *Service {
	return &Service{taskRepo: taskRepo, noteRepo: noteRepo, groupRepo: groupRepo}
}

func (s *Service) Create(task *model.Task) error {
	return s.taskRepo.Create(task)
}

func (s *Service) FindByID(id string) (*model.Task, error) {
	return s.taskRepo.FindByID(id)
}

func (s *Service) FindAll() ([]model.Task, error) {
	return s.taskRepo.FindAll()
}

func (s *Service) FindNoteByTaskID(taskID string) (*model.NoteRecord, error) {
	return s.noteRepo.FindByTaskID(taskID)
}

func (s *Service) Delete(id string) error {
	s.noteRepo.DeleteByTaskID(id)
	return s.taskRepo.Delete(id)
}

func (s *Service) UpdateName(id, name string) error {
	return s.taskRepo.UpdateName(id, name)
}

func (s *Service) FindAllGroups() ([]model.Group, error) {
	return s.groupRepo.FindAll()
}
