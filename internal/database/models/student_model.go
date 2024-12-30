package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

type Students struct {
	BaseModel
	ID       uint    `gorm:"primaryKey;autoIncrement"`
	Name     string  `gorm:"not null"`
	Lastname string  `gorm:"not null"`
	Email    *string `gorm:"uniqueIndex"`
	Number   *string `gorm:"uniqueIndex"`
	Gender   *Gender `gorm:"check:gender IN ('Male','Female')"`
	Birthday *string

	CoursesId []uint `gorm:"-"`

	Courses []Courses `gorm:"many2many:students_courses;joinForeignKey:StudentsID;joinReferences:CoursesID"`
}

type StudentWithEntrance struct {
	Students
	Entrance time.Time `json:"entrance"`
}

func (s *Students) BeforeSave(tx *gorm.DB) (err error) {
	if s.Email == nil && s.Number == nil {
		return fmt.Errorf("al menos uno de los campos email o number debe estar presente")
	}
	return
}

func (s *Students) BeforeDelete(tx *gorm.DB) (err error) {
	err = tx.Model(s).Association("Courses").Clear()
	if err != nil {
		return fmt.Errorf("error al eliminar relaciones con cursos: %w", err)
	}
	return
}
