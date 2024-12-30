package courses

type CreateOrUpdateCourse struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
}

type QueryGetCourses struct {
	Page      int `form:"page"`
	Limit     int `form:"limit"`
	StudentId int `form:"studentId"`
}
