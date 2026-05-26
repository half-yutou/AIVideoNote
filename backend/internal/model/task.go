package model

import (
	"time"
)

type Task struct {
	ID           string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	VideoURL     string     `gorm:"not null" json:"video_url"`
	Platform     string     `gorm:"not null;size:20" json:"platform"`
	VideoID      string     `gorm:"size:100" json:"video_id"`
	Status       string     `gorm:"not null;size:20;default:PENDING" json:"status"`
	Quality      string     `gorm:"size:10;default:medium" json:"quality"`
	ModelName    string     `gorm:"size:50" json:"model_name"`
	ProviderID   string     `gorm:"size:36" json:"provider_id"`
	Format       string     `gorm:"type:text;default:'[]'" json:"format"`
	Style        string     `gorm:"size:50" json:"style"`
	Extras       string     `gorm:"type:text" json:"extras"`
	Name         string     `gorm:"size:200" json:"name"`
	GroupID      string     `gorm:"size:36" json:"group_id"`
	Link         bool       `gorm:"default:false" json:"link"`
	ErrorMessage string     `gorm:"type:text" json:"error_message"`
	RetryCount   int        `gorm:"default:0" json:"retry_count"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CompletedAt  *time.Time `json:"completed_at"`
}

func (Task) TableName() string {
	return "tasks"
}

type TaskStatus = string

const (
	TaskStatusPending              TaskStatus = "PENDING"
	TaskStatusParsing              TaskStatus = "PARSING"
	TaskStatusDownloading          TaskStatus = "DOWNLOADING"
	TaskStatusTranscribing         TaskStatus = "TRANSCRIBING"
	TaskStatusGenerating           TaskStatus = "GENERATING"
	TaskStatusPostProcessing       TaskStatus = "POST_PROCESSING"
	TaskStatusIndexing             TaskStatus = "INDEXING"
	TaskStatusSuccess              TaskStatus = "SUCCESS"
	TaskStatusFailed               TaskStatus = "FAILED"
	TaskStatusParsingFailed        TaskStatus = "PARSING_FAILED"
	TaskStatusDownloadingFailed    TaskStatus = "DOWNLOADING_FAILED"
	TaskStatusTranscribingFailed   TaskStatus = "TRANSCRIBING_FAILED"
	TaskStatusGeneratingFailed     TaskStatus = "GENERATING_FAILED"
	TaskStatusPostProcessingFailed TaskStatus = "POST_PROCESSING_FAILED"
)
