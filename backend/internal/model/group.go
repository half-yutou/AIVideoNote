package model

import "time"

type Group struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name      string    `gorm:"not null;size:100" json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (Group) TableName() string {
	return "groups"
}
