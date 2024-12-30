package courses

import (
	"github.com/juseph-q/SchoolPr/internal/database/models"
	"github.com/juseph-q/SchoolPr/internal/shared/handler"
)

type CourseResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CoursesResponse struct {
	Courses []course             `json:"courses"`
	Meta    handler.MetadataPage `json:"metadata"`
}

type course struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	TotalStudents int64  `json:"totalStudents"`
}

func NewCourseResponse(courseDb *models.Courses) *CourseResponse {

	return &CourseResponse{
		ID:   courseDb.ID,
		Name: courseDb.Name,
	}
}

func NewCoursesResponse(courseDb []*models.CourseWithStudentCount, metadata handler.MetadataPage) *CoursesResponse {
	var courses []course

	for _, course := range courseDb {
		courses = append(courses, *newCourse(course))
	}

	return &CoursesResponse{
		Courses: courses,
		Meta:    metadata,
	}
}

func newCourse(courseDb *models.CourseWithStudentCount) *course {
	return &course{
		ID:            courseDb.ID,
		Name:          courseDb.Name,
		TotalStudents: courseDb.StudentCount,
	}
}
