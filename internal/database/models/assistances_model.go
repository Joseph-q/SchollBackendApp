package models

import "time"

type Assistance struct {
	StudentID uint      `gorm:"not null;primaryKey"`
	CourseID  uint      `gorm:"not null;primaryKey"`
	Date      time.Time `gorm:"not null;primaryKey;type:date"`
	Time      string    `gorm:"not null;type:time"`

	// Relaciones
	Students Students `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE"`
	Courses  Courses  `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE"`
}
