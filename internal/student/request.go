package student

import "github.com/juseph-q/SchoolPr/internal/database/models"

type QueryGetStudentById struct {
	CourseId *string `form:"courseId"`
}

type QuerySearchStudent struct {
	Query  string  `form:"query"`
	Limit  int     `form:"limit"`
	Name   *string `form:"name"`
	Email  *string `form:"email"`
	Number *string `form:"number"`
}

type QueryGetStudents struct {
	Limit    int     `form:"limit" json:"limit"`
	Page     int     `form:"page" json:"page"`
	Order    string  `form:"orderBy" json:"orderBy"`
	CourseId *string `form:"courseId" json:"courseId"`
}

type CreateStudent struct {
	Name      string         `json:"name" binding:"required,min=2,max=50"`
	Lastname  string         `json:"lastname" binding:"required,min=2,max=50"`
	Email     *string        `json:"email,omitempty" binding:"omitempty,required,email"`
	Number    *string        `json:"number,omitempty" binding:"omitempty,required"`
	Gender    *models.Gender `json:"gender,omitempty" binding:"omitempty,gender"`
	Birthday  *string        `json:"birthday,omitempty" binding:"omitempty,dateformat"`
	CoursesID []uint         `json:"coursesId,omitempty" binding:"omitempty,dive,gt=0"`
}

type UpdateStudent struct {
	Name      *string        `json:"name,omitempty"`
	Lastname  *string        `json:"lastname,omitempty"`
	Email     *string        `json:"email,omitempty" binding:"omitempty,required,email"`
	Number    *string        `json:"number,omitempty" binding:"omitempty,required"`
	Gender    *models.Gender `json:"gender,omitempty" binding:"omitempty,gender"`
	Birthday  *string        `json:"birthday,omitempty" binding:"omitempty,dateformat"`
	CoursesID []uint         `json:"coursesId,omitempty" binding:"omitempty,dive,gt=0"`
}

func NewStudentUpdateDb(studentTU *UpdateStudent, studentDb *models.Students) *models.Students {
	if studentTU.Name != nil {
		studentDb.Name = *studentTU.Name
	}

	// Lastname
	if studentTU.Lastname != nil {
		studentDb.Lastname = *studentTU.Lastname
	}

	// Email
	if studentTU.Email != nil {
		studentDb.Email = studentTU.Email
	}

	// Number
	if studentTU.Number != nil {
		studentDb.Number = studentTU.Number
	}

	// Gender
	if studentTU.Gender != nil {
		studentDb.Gender = studentTU.Gender
	}

	// Birthday
	if studentTU.Birthday != nil {
		studentDb.Birthday = studentTU.Birthday
	}

	if studentTU.Birthday != nil {
		studentDb.Birthday = studentTU.Birthday
	}

	if len(studentTU.CoursesID) > 0 {
		studentDb.CoursesId = studentTU.CoursesID
	}

	return studentDb
}
