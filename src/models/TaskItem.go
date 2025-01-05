package models

import "time"

type TaskItem struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Priority    string    `gorm:"not null"`
	Description string    `gorm:"default:NULL"`
	DueDate     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status      string    `gorm:"status"`
	Category    string    `gorm:"default:NULL"`
	Project     string    `gorm:"default:NULL"`
}