package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Completed bool           `json:"completed" gorm:"default:false"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}
