package model

import (
	"time"
)

type NoteRecord struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TaskID     string    `gorm:"not null;size:36;index" json:"task_id"`
	Markdown   string    `gorm:"not null;type:text" json:"markdown"`
	Transcript string    `gorm:"type:text" json:"transcript"`
	AudioMeta  string    `gorm:"type:text" json:"audio_meta"`
	Version    int       `gorm:"default:1" json:"version"`
	CreatedAt  time.Time `json:"created_at"`
}

func (NoteRecord) TableName() string {
	return "notes"
}
