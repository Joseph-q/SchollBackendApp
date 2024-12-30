package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
