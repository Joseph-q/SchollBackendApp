package models

import "time"

type AssistanceHistorial struct {
	CourseID uint      `gorm:"not null;primaryKey"`
	Date     time.Time `gorm:"not null;primaryKey;type:date"`
	Time     string    `gorm:"not null;type:time"`

	Courses Courses `gorm:"foreignKey:CourseID"`
}
