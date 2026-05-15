package models

import "time"

type Task struct {
	ID          uint      `json:"id" gorm:"primary_key;auto_increment"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Priority    int       `json:"priority" gorm:"default:1"`
	Status      string    `json:"status" gorm:"default:'pending'"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
