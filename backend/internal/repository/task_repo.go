package repository

import (
	"strings"

	"github.com/google/uuid"
	"github.com/aivideonote/backend/internal/database"
	"github.com/aivideonote/backend/internal/model"
)

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) Create(task *model.Task) error {
	if task.ID == "" {
		task.ID = uuid.New().String()
	}
	task.Status = model.TaskStatusPending
	return database.DB.Create(task).Error
}

func (r *TaskRepository) FindByID(id string) (*model.Task, error) {
	var task model.Task
	err := database.DB.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) FindAll() ([]model.Task, error) {
	var tasks []model.Task
	err := database.DB.Order("created_at DESC").Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) FindByVideoID(videoID, platform string) (*model.Task, error) {
	var task model.Task
	err := database.DB.Where("video_id = ? AND platform = ?", videoID, platform).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) UpdateStatus(id string, status model.TaskStatus, errMsg string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if errMsg != "" {
		updates["error_message"] = errMsg
	}
	if status == model.TaskStatusSuccess || status == model.TaskStatusFailed || strings.HasSuffix(string(status), "_FAILED") {
		updates["completed_at"] = database.DB.Raw("NOW()")
	}
	return database.DB.Model(&model.Task{}).Where("id = ?", id).Updates(updates).Error
}

func (r *TaskRepository) UpdateVideoID(id, videoID string) error {
	return database.DB.Model(&model.Task{}).Where("id = ?", id).Update("video_id", videoID).Error
}

func (r *TaskRepository) Delete(id string) error {
	return database.DB.Delete(&model.Task{}, "id = ?", id).Error
}

func (r *TaskRepository) UpdateName(id, name string) error {
	return database.DB.Model(&model.Task{}).Where("id = ?", id).Update("name", name).Error
}
