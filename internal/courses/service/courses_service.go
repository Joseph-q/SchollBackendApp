package coursesService

import (
	"context"
	"errors"
	"strings"

	"github.com/juseph-q/SchoolPr/internal/courses"
	"github.com/juseph-q/SchoolPr/internal/database"
	"github.com/juseph-q/SchoolPr/internal/database/models"
	"gorm.io/gorm"
)

type CourseService struct {
	db *gorm.DB
}

func NewCourseService(db *gorm.DB) *CourseService {
	return &CourseService{
		db: db,
	}
}

func (s *CourseService) FindCourses(queryReq *courses.QueryGetCourses, ctx context.Context) (courseDb []*models.CourseWithStudentCount, coursesTotal int64, err error) {
	var coursesWithCountStudent []*models.CourseWithStudentCount
	var totalCourses int64

	db := database.FromContext(ctx, s.db)
	chain := db.WithContext(ctx)

	offset := (queryReq.Page - 1) * queryReq.Limit

	// Perform a query to obtain the classrooms and student count in a single query
	if err := chain.Model(&models.Courses{}).
		Select("courses.id, courses.name, COUNT(students_courses.students_id) AS student_count").
		Joins("LEFT JOIN students_courses ON students_courses.courses_id = courses.id").
		Joins("LEFT JOIN students ON students_courses.students_id = students.id").
		Group("courses.id").
		Offset(offset).
		Limit(queryReq.Limit).
		Scan(&coursesWithCountStudent).Error; err != nil {
		return nil, 0, err
	}

	if err := chain.Model(&models.Courses{}).Count(&totalCourses).Error; err != nil {
		return nil, 0, err
	}

	return coursesWithCountStudent, totalCourses, nil
}

func (s *CourseService) FindCourseById(id uint, ctx context.Context) (*models.Courses, error) {
	db := database.FromContext(ctx, s.db)
	var courseDb models.Courses

	if err := db.WithContext(ctx).Model(&models.Courses{}).First(&courseDb, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, courses.CourseNotFound
		}
		return nil, err
	}

	return &courseDb, nil
}

func (s *CourseService) CreateCourse(courseTC *courses.CreateOrUpdateCourse, ctx context.Context) (*models.Courses, error) {
	var courseDb = models.Courses{
		Name: courseTC.Name,
	}

	db := database.FromContext(ctx, s.db)
	if err := db.Create(&courseDb).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, courses.CourseNameIsRegister
		}
		return nil, err
	}

	return &courseDb, nil

}

func (s *CourseService) UpdateCourse(ctx context.Context, data *models.Courses) error {
	db := database.FromContext(ctx, s.db)

	if err := db.WithContext(ctx).Model(&data).Updates(data).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return courses.CourseNameIsRegister
		}
		return err
	}

	return nil
}
