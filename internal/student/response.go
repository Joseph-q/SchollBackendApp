package student

import (
	"time"

	"github.com/juseph-q/SchoolPr/internal/database/models"
	"github.com/juseph-q/SchoolPr/internal/shared/handler"
)

type StudentResponse struct {
	Id        uint           `json:"id"`
	Name      string         `json:"name"`
	Lastname  string         `json:"lastname"`
	Email     *string        `json:"email"`
	Number    *string        `json:"number"`
	Gender    *models.Gender `json:"gender"`
	Birthday  *string        `json:"birthday"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Courses   []Courses      `json:"courses"`
}

type Courses struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type StudentsResponse struct {
	Student []StudentsFromStudents `json:"data"`
	Meta    handler.MetadataPage   `json:"metadata"`
	Query   QueryGetStudents       `json:"queryParams"`
}

type StudentsFromStudents struct {
	Id        uint           `json:"id"`
	Name      string         `json:"name"`
	Lastname  string         `json:"lastname"`
	Email     *string        `json:"email"`
	Number    *string        `json:"number"`
	Gender    *models.Gender `json:"gender"`
	Birthday  *string        `json:"birthday"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type StudentSearchResponse struct {
	Hits               []studentFromSearch `json:"hits"`
	EstimatedTotalHits int64               `json:"estimatedTotalHits"`
	Limit              int                 `json:"limit"`
	Query              string              `json:"query"`
}

type studentFromSearch struct {
	Id       uint    `json:"id"`
	Name     string  `json:"name"`
	Lastname string  `json:"lastname"`
	Email    *string `json:"email"`
	Number   *string `json:"number"`
}

func NewStudentResponse(student *models.Students) *StudentResponse {

	var courses []Courses

	for _, course := range student.Courses {
		courses = append(courses, Courses{
			Id:   course.ID,
			Name: course.Name,
		})
	}

	return &StudentResponse{
		Id:        student.ID,
		Name:      student.Name,
		Lastname:  student.Lastname,
		Email:     student.Email,
		Number:    student.Number,
		Gender:    student.Gender,
		Birthday:  student.Birthday,
		CreatedAt: student.CreatedAt,
		UpdatedAt: student.UpdatedAt,
		Courses:   courses,
	}
}

func NewStudentsResponse(students []models.Students, meta handler.MetadataPage, query QueryGetStudents) *StudentsResponse {
	var responses []StudentsFromStudents

	for _, student := range students {
		responses = append(responses, StudentsFromStudents{
			Id:        student.ID,
			Name:      student.Name,
			Lastname:  student.Lastname,
			Email:     student.Email,
			Number:    student.Number,
			Gender:    student.Gender,
			Birthday:  student.Birthday,
			CreatedAt: student.CreatedAt,
			UpdatedAt: student.UpdatedAt,
		})
	}

	return &StudentsResponse{
		Student: responses,
		Meta:    meta,
		Query:   query,
	}
}

func NewStudentSearchReponse(students []models.Students, totalHits *int64, limit int, query string) *StudentSearchResponse {
	var response []studentFromSearch
	for _, student := range students {
		response = append(response, studentFromSearch{
			Id:       student.ID,
			Name:     student.Name,
			Email:    student.Email,
			Number:   student.Number,
			Lastname: student.Lastname,
		})
	}

	return &StudentSearchResponse{
		Hits:               response,
		EstimatedTotalHits: *totalHits,
		Limit:              limit,
		Query:              query,
	}
}
