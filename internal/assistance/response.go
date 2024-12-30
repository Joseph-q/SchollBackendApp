package assistance

import (
	"time"

	"github.com/juseph-q/SchoolPr/internal/database/models"
	"github.com/juseph-q/SchoolPr/internal/shared/handler"
)

type ResponseCourseAssistances struct {
	Data DataResponseCourseAssistance `json:"data"`
	Meta *handler.MetadataPage        `json:"metadata"`
}
type DataResponseCourseAssistance struct {
	Date     string            `json:"date"`
	Students []StudentEntrance `json:"students"`
}

type StudentEntrance struct {
	Id       uint    `json:"id"`
	Name     string  `json:"name"`
	Lastname string  `json:"lastname"`
	Email    *string `json:"email"`
	Number   *string `json:"number"`
	Entrance string  `json:"entrance"`
}

type ResponseHistorialAssistance struct {
	Data []DataHistorial `json:"data"`
}

type DataHistorial struct {
	CourseId  uint              `json:"courseId"`
	Assitance []AssistanceCount `json:"assistances"`
}
type AssistanceCount struct {
	Date            string `json:"date"`
	TotalAssistance int64  `json:"total"`
}

type ResponseAssistancesStudent struct {
	StudentId   uint                `json:"studentId"`
	Assistances []courseAssistances `json:"courses"`
}

type assistance struct {
	Date     time.Time `json:"date"`
	Entrance string    `json:"entrance"`
}

type courseAssistances struct {
	CourseId    uint         `json:"id"`
	CourseName  string       `json:"name"`
	Assistances []assistance `json:"assistances"`
}

type HistorialAssist struct {
	CourseId  uint   `json:"courseId"`
	StudentId uint   `json:"studentId"`
	Date      string `json:"date"`
	Total     int64  `json:"total"`
}
type ResponseAssistanceSummary struct {
	Data []AssistanceSummary   `json:"data"`
	Meta *handler.MetadataPage `json:"metadata"`
}

type AssistanceSummary struct {
	Date  string `json:"date"`
	Total int    `json:"total"`
}

type ResponseCourseHistorial struct {
	Date string          `json:"date"`
	Data []DataHistorial `json:"course"`
}

type CourseHistorial struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	TotalAssistance int64  `json:"total"`
}

func NewCourseResponseAssistance(data []models.Assistance, date string, metadata handler.MetadataPage) ResponseCourseAssistances {
	var studentEntrance []StudentEntrance
	for _, assitance := range data {
		studentEntrance = append(studentEntrance, StudentEntrance{
			Id:       assitance.Students.ID,
			Name:     assitance.Students.Name,
			Lastname: assitance.Students.Lastname,
			Email:    assitance.Students.Email,
			Number:   assitance.Students.Number,
			Entrance: assitance.Time,
		})
	}

	return ResponseCourseAssistances{
		Data: DataResponseCourseAssistance{
			Date:     date,
			Students: studentEntrance,
		},
		Meta: &metadata,
	}
}

func NewResponseAssistancesStudent(studentId uint, assitances []*models.Assistance) *ResponseAssistancesStudent {

	courseMap := make(map[uint]*courseAssistances)

	for _, assist := range assitances {
		if _, exist := courseMap[assist.CourseID]; !exist {
			courseMap[assist.CourseID] = &courseAssistances{
				CourseId:   assist.CourseID,
				CourseName: assist.Courses.Name,
			}
		}
		courseMap[assist.CourseID].Assistances = append(courseMap[assist.CourseID].Assistances, assistance{
			Date:     assist.Date,
			Entrance: assist.Time,
		})
	}
	dataCourse := make([]courseAssistances, 0, len(courseMap))
	for _, course := range courseMap {
		dataCourse = append(dataCourse, *course)
	}

	return &ResponseAssistancesStudent{
		StudentId:   studentId,
		Assistances: dataCourse,
	}

}

func NewResponseHistorial(input []HistorialAssist) ResponseHistorialAssistance {
	// Map to group assistance by courseId
	courseMap := make(map[uint][]AssistanceCount)

	for _, record := range input {
		assistance := AssistanceCount{
			Date:            record.Date,
			TotalAssistance: record.Total,
		}
		courseMap[record.CourseId] = append(courseMap[record.CourseId], assistance)
	}

	// Build the final response
	response := ResponseHistorialAssistance{}
	for courseId, assistances := range courseMap {
		courseIdCopy := courseId // Create a copy to avoid pointer issues
		response.Data = append(response.Data, DataHistorial{
			CourseId:  courseIdCopy,
			Assitance: assistances,
		})
	}

	return response
}

func NewResponseAssistancesSumary(data []AssistanceSummary, metadata handler.MetadataPage) *ResponseAssistanceSummary {
	return &ResponseAssistanceSummary{
		Data: data,
		Meta: &metadata,
	}
}
