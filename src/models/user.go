package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"default:NULL"`
	Phone     string    `gorm:"default:NULL"`
	DOB       time.Time // Date of birth of the user
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"` // When the user was created in the system
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"` // When the user was last updated in the system
}