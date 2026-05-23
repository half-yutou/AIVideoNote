package model

import (
	"time"
)

type LLMProvider struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string    `gorm:"not null;size:100" json:"name"`
	ApiKey    string    `gorm:"not null;type:text" json:"-"`
	BaseURL   string    `gorm:"not null;type:text" json:"base_url"`
	Type      string    `gorm:"not null;size:20" json:"type"`
	Logo      string    `gorm:"type:text" json:"logo"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (LLMProvider) TableName() string {
	return "providers"
}
