package export

import (
	"os"
	"path/filepath"

	"github.com/aivideonote/backend/internal/config"
	"github.com/aivideonote/backend/internal/repository"
)

type Service struct {
	cfg      *config.Config
	noteRepo *repository.NoteRepository
	taskRepo *repository.TaskRepository
}

func NewService(cfg *config.Config, noteRepo *repository.NoteRepository, taskRepo *repository.TaskRepository) *Service {
	return &Service{cfg: cfg, noteRepo: noteRepo, taskRepo: taskRepo}
}

func (s *Service) ExportMarkdown(taskID string) (string, error) {
	note, err := s.noteRepo.FindByTaskID(taskID)
	if err != nil {
		return "", err
	}
	return note.Markdown, nil
}

func (s *Service) ExportFile(taskID string) (string, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return "", err
	}
	mdPath := filepath.Join(s.cfg.Storage.DataDir, task.VideoID, taskID+".md")
	if _, err := os.Stat(mdPath); os.IsNotExist(err) {
		return "", err
	}
	return mdPath, nil
}
