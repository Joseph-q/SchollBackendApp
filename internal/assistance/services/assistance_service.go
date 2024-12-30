package assistanceServices

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/juseph-q/SchoolPr/internal/assistance"
	"github.com/juseph-q/SchoolPr/internal/database"
	"github.com/juseph-q/SchoolPr/internal/database/models"
	"github.com/juseph-q/SchoolPr/internal/shared/handler"
)

type AssistanceService struct {
	db *gorm.DB
}

func NewAssitanceService(db *gorm.DB) *AssistanceService {
	return &AssistanceService{
		db: db,
	}
}

func (s *AssistanceService) StudentsAssitanceByCourseId(course *models.Courses, queryReq *assistance.QueryParamsByCourseId, ctx context.Context) ([]models.Assistance, int64, error) {
	db := database.FromContext(ctx, s.db)

	var resultStudent []struct {
		ID     uint    `json:"id"`
		Name   string  `json:"name"`
		Email  *string `json:"email"`
		Number *string `json:"number"`
		Time   string  `json:"time"`
	}
	var totalStudentsAssitances int64

	chain := db.WithContext(ctx).
		Table("assistances").
		Joins("JOIN students ON assistances.student_id = students.id").
		Select("students.id, students.email, students.number, students.name, assistances.time")

	if queryReq != nil {
		if queryReq.Date != "" {
			chain = chain.Where("date = ?", queryReq.Date)
		}
		if queryReq.StudentId > 0 {
			chain = chain.Where("student_id = ?", queryReq.StudentId)
		}

	}

	offset := (queryReq.Page - 1) * queryReq.Limit

	chain = chain.Offset(offset).
		Limit(queryReq.Limit).
		Where(&models.Assistance{CourseID: course.ID}).
		Count(&totalStudentsAssitances).
		Scan(&resultStudent)

	if err := chain.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, assistance.NotFountRegister
		}
		return nil, 0, err
	}
	var assitance []models.Assistance

	for _, student := range resultStudent {
		assitance = append(assitance, models.Assistance{
			Students: models.Students{
				ID:     student.ID,
				Name:   student.Name,
				Email:  student.Email,
				Number: student.Number,
			},
			Time: student.Time,
		})
	}

	return assitance, totalStudentsAssitances, nil
}

func (s *AssistanceService) FindStudentAssistance(ctx context.Context, id int, queryReq *assistance.QueryParamsAssistanceStudentById) ([]*models.Assistance, error) {

	db := database.FromContext(ctx, s.db)

	var result []*models.Assistance

	chain := db.WithContext(ctx).Model(&models.Assistance{}).
		Preload("Courses").
		Where("assistances.student_id = ?", id)
	var orderClause = "date DESC"

	if queryReq != nil {

		if queryReq.CourseId > 0 {
			chain = chain.Where("assistances.course_id = ?", queryReq.CourseId)
		}

		if queryReq.Date != "" {
			chain = chain.Where("assistances.date = ?", queryReq.Date)
		}

		switch queryReq.OrderBy {
		case "dateAsc":
			orderClause = "date ASC"
		case "dateDesc":
			orderClause = "date DESC"
		}

	}
	chain = chain.Order(orderClause)

	if err := chain.Find(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, assistance.NotFountRegister
		}
		return nil, handler.ServerError
	}

	return result, nil
}

type AssistanceWithCourse struct {
	CourseID uint   `json:"id"`   // ID of the course
	Name     string `json:"name"` // Name of the course
}

func (s *AssistanceService) FindStudentCoursesAssisted(ctx context.Context, id int) ([]*AssistanceWithCourse, error) {
	db := database.FromContext(ctx, s.db)

	var results []*AssistanceWithCourse

	err := db.WithContext(ctx).
		Table("assistances AS a").                     // Use the "assistances" table with alias 'a'
		Select("a.course_id, c.name").                 // Select the course_id and course name
		Joins("JOIN courses c ON a.course_id = c.id"). // Join with the "courses" table using the course ID
		Where("a.student_id = ?", id).                 // Filter by student_id = 1
		Group("a.course_id, c.name").                  // Group by course_id and course name to ensure unique results
		Scan(&results).                                // Scan the results into the AssistanceWithCourse struct
		Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *AssistanceService) HistorialAssistances(ctx context.Context, queryParams *assistance.QueryParamsHistorialAssitance) ([]assistance.HistorialAssist, error) {

	db := database.FromContext(ctx, s.db)

	chain := db.WithContext(ctx).Model(&models.Assistance{}).Select("course_id,date, COUNT(*) AS total").Group("course_id, date")

	var resultsQuery []assistance.HistorialAssist

	if queryParams != nil {

		if queryParams.CourseId != "" {
			fmt.Println(queryParams.CourseId)
			chain = chain.Where("course_id = ?", queryParams.CourseId)
		}
		if queryParams.StudentId != "" {
			chain = chain.Where("student_id = ?", queryParams.StudentId)
		}

		if queryParams.Date != "" {
			chain = chain.Where("date = ?", queryParams.Date)
		}
		if queryParams.StartDate != "" && queryParams.EndDate != "" {
			if queryParams.StartDate == queryParams.EndDate {
				chain = chain.Where("date = ?", queryParams.StartDate)
			} else {
				chain = chain.Where("date BETWEEN ? AND ?", queryParams.StartDate, queryParams.EndDate)
			}
		}

	}

	if err := chain.Order("date DESC").Scan(&resultsQuery).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, assistance.NotFountRegister
		}

		return nil, handler.ServerError
	}

	return resultsQuery, nil
}

func (s *AssistanceService) HistorialAssistancesSumary(ctx context.Context, queryParams *assistance.QueryParamsHistorialAssitanceSummary) ([]assistance.AssistanceSummary, error) {
	db := database.FromContext(ctx, s.db)
	chain := db.Model(models.Assistance{}).Select("date", "COUNT(student_id) as total").Group("date")

	if queryParams != nil {
		if queryParams.EndDate != "" && queryParams.StartDate != "" {
			chain = chain.Where("date BETWEEN ? AND ?", queryParams.StartDate, queryParams.EndDate)
		}
	}

	var assistSumary []assistance.AssistanceSummary
	if err := chain.Scan(&assistSumary).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, assistance.NotFountRegister
		}
		return nil, err
	}

	return assistSumary, nil
}
