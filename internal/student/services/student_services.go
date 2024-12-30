package studentService

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/juseph-q/SchoolPr/internal/database"
	"github.com/juseph-q/SchoolPr/internal/database/models"
	"github.com/juseph-q/SchoolPr/internal/shared/handler"
	"github.com/juseph-q/SchoolPr/internal/student"
	"gorm.io/gorm"
)

type StudentService struct {
	db *gorm.DB
}

func NewStudentService(db *gorm.DB) *StudentService {
	return &StudentService{
		db: db,
	}
}

func (s *StudentService) FindStudentById(id int, queryReq *student.QueryGetStudentById, ctx context.Context) (*models.Students, error) {
	var studentDB models.Students

	db := database.FromContext(ctx, s.db)

	tx := db.WithContext(ctx).Model(&models.Students{})

	if queryReq != nil {
		if queryReq.CourseId != nil && *queryReq.CourseId == "all" {
			tx = tx.Preload("Courses") //Change to join if want speed
		} else if queryReq.CourseId != nil {
			tx = tx.Preload("Courses", "courses.id = ?", *queryReq.CourseId) //Change to join if want speed
		}
	}

	if err := tx.Where("students.id = ?", id).First(&studentDB).Error; err != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, student.UserNotFound
		}
		return nil, err
	}

	return &studentDB, nil
}

func (s *StudentService) FindStudents(queryReq student.QueryGetStudents, ctx context.Context) ([]models.Students, int64, error) {
	var students []models.Students
	db := database.FromContext(ctx, s.db)

	var orderClause string

	offset := (queryReq.Page - 1) * queryReq.Limit
	var totalStudents int64

	switch queryReq.Order {
	case "nameAsc":
		orderClause = "name ASC"
	case "nameDesc":
		orderClause = "name DESC"
	case "createdAsc":
		orderClause = "created_at ASC"
	case "createdDesc":
		orderClause = "created_at DESC"
	case "updatedAsc":
		orderClause = "updated_at ASC"
	case "updatedDesc":
		orderClause = "updated_at DESC"
	default:
		orderClause = "created_at DESC"
	}

	tx := db.WithContext(ctx).
		Model(&models.Students{})

	if queryReq.CourseId != nil {
		if *queryReq.CourseId == "null" {
			tx = tx.Joins("LEFT JOIN students_courses ON students_courses.students_id = students.id").
				Where("students_courses.courses_id IS NULL")
		} else {
			tx = tx.Joins("JOIN students_courses ON students_courses.students_id = students.id").
				Where("students_courses.courses_id = ?", *queryReq.CourseId)
		}
	}

	tx.Count(&totalStudents)

	if err := tx.Offset(offset).
		Limit(queryReq.Limit).
		Order(orderClause).Find(&students).Error; err != nil {
		return nil, 0, handler.ServerError
	}

	return students, totalStudents, nil
}

func (s *StudentService) CreateStudent(studentTC student.CreateStudent, ctx context.Context) (*models.Students, error) {
	db := database.FromContext(ctx, s.db)

	tx := db.WithContext(ctx).Begin()

	if err := tx.Error; err != nil {
		return nil, err
	}

	newStudent := models.Students{
		Name:     studentTC.Name,
		Lastname: studentTC.Lastname,
		Email:    studentTC.Email,
		Number:   studentTC.Number,
		Gender:   studentTC.Gender,
		Birthday: studentTC.Birthday,
	}

	// Guardar el estudiante en la base de datos dentro de la transacción
	if err := tx.Create(&newStudent).Error; err != nil {
		tx.Rollback() // Deshacer la transacción en caso de error
		fmt.Println(err.Error())
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return nil, student.EmailOrNumberAlreadyRegister
		}
		return nil, err
	}

	if len(studentTC.CoursesID) > 0 {
		var courses []models.Courses
		// Buscar los cursos por sus IDs dentro de la transacción
		if err := tx.Where("id IN (?)", studentTC.CoursesID).Find(&courses).Error; err != nil {
			tx.Rollback() // Deshacer la transacción en caso de error
			return nil, err
		}

		if err := tx.Model(&newStudent).Association("Courses").Append(courses); err != nil {
			if errors.Is(err, gorm.ErrInvalidData) || errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback() // Deshacer la transacción en caso de error crítico de base de datos
				return nil, fmt.Errorf("error al asociar los cursos: %w", err)
			}
			fmt.Println("Error no crítico al asociar los cursos:", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // Deshacer la transacción en caso de error
		return nil, fmt.Errorf("error al confirmar la transacción")
	}

	return &newStudent, nil
}

func (s *StudentService) UpdateStudent(studentTU *models.Students, ctx context.Context) error {
	db := database.FromContext(ctx, s.db)

	fields := make(map[string]interface{})

	tx := db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	if len(studentTU.CoursesId) > 0 {
		deleteQuery := `DELETE FROM students_courses WHERE students_id = ?;`
		if err := tx.Exec(deleteQuery, studentTU.ID).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Construcción de la consulta INSERT
		insertQuery := `INSERT INTO students_courses (students_id, courses_id) VALUES `
		var values []interface{}
		var sb strings.Builder
		sb.Grow(len(studentTU.CoursesId) * 2) // Pre-alocar memoria para la cadena

		for _, course := range studentTU.CoursesId {
			sb.WriteString("(?, ?),")
			values = append(values, studentTU.ID, course)
		}

		// Eliminar la última coma extra
		insertQuery += sb.String()[:sb.Len()-1]

		if err := tx.Exec(insertQuery, values...).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if studentTU.Name != "" {
		fields["name"] = studentTU.Name
	}

	if studentTU.Lastname != "" {
		fields["lastname"] = studentTU.Lastname
	}

	if studentTU.Email != nil && *studentTU.Email != "" {
		fields["email"] = *studentTU.Email
	}

	if studentTU.Number != nil && *studentTU.Number != "" {
		fields["number"] = *studentTU.Number
	}

	if studentTU.Gender != nil {
		fields["gender"] = *studentTU.Gender
	}

	if studentTU.Birthday != nil && *studentTU.Birthday != "" {
		fields["birthday"] = *studentTU.Birthday
	}

	if len(fields) > 0 {
		if err := tx.Model(&studentTU).Updates(fields).Error; err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				tx.Rollback()
				return student.EmailOrNumberAlreadyRegister
			}
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *StudentService) DeleteStudent(studentTD *models.Students, ctx context.Context) error {
	db := database.FromContext(ctx, s.db)
	chain := db.WithContext(ctx).Model(&models.Students{}).Delete(studentTD)

	if err := chain.Error; err != nil {
		return chain.Error
	}

	if chain.RowsAffected == 0 {
		return student.UserNotFound
	}

	return nil
}

func (s *StudentService) SearchStudent(queryReq *student.QuerySearchStudent, ctx context.Context) ([]models.Students, *int64, error) {
	var data []models.Students
	var total int64
	db := database.FromContext(ctx, s.db)

	chain := db.Model(&models.Students{}).Select("id", "name", "lastname", "email", "number").Where("deleted_at IS NULL")

	if queryReq == nil {
		return data, &total, nil
	}

	if queryReq.Email != nil {
		chain = chain.Where("LOWER(email) LIKE ?", "%"+*queryReq.Email+"%")

	} else if queryReq.Number != nil {
		chain = chain.Where("number LIKE ?", *queryReq.Number+"%")
	} else if queryReq.Name != nil {
		chain = chain.Where(`
		LOWER(REPLACE(REPLACE(REPLACE(name, '.', ''), ',', ''), '-', '')) LIKE ?
		OR LOWER(REPLACE(REPLACE(REPLACE(lastname, '.', ''), ',', ''), '-', '')) LIKE ?`,
			strings.ToLower(*queryReq.Name)+"%", strings.ToLower(*queryReq.Name)+"%")
	}

	if queryReq.Limit != 0 {
		chain = chain.Limit(queryReq.Limit)
	} else {
		chain = chain.Limit(10)
	}

	if err := chain.Find(&data).Count(&total).Error; err != nil {
		return nil, nil, err
	}

	return data, &total, nil
}
