package assistance

type QueryParamsByCourseId struct {
	Page      int    `form:"page,omitempty" binding:"omitempty"`
	Limit     int    `form:"limit,omitempty"  binding:"omitempty"`
	Date      string `form:"date,omitempty" binding:"omitempty,dateformat"`
	StudentId uint   `form:"studentId,omitempty" binding:"omitempty"`
}

type QueryParamsAssistanceStudentById struct {
	CourseId uint   `form:"courseId"`
	OrderBy  string `form:"orderBy"`
	Date     string `form:"date,omitempty" binding:"omitempty,dateformat"`
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
}

type QueryParamsHistorialAssitance struct {
	CourseId  string `form:"courseId"`
	StudentId string `form:"studentId"`
	Date      string `form:"date"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

type QueryParamsHistorialAssitanceSummary struct {
	CourseId  string `form:"courseId"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
}

type QueryParamsAssitanceSumary struct {
	Date []string `form:"date"`
}
