package model

import "time"

type PlatformCookie struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Platform  string    `gorm:"not null;uniqueIndex;size:20" json:"platform"`
	Content   string    `gorm:"not null;type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PlatformCookie) TableName() string {
	return "platform_cookies"
}
