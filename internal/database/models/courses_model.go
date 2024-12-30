package models

type Courses struct {
	BaseModel
	ID       uint       `gorm:"primaryKey;autoIncrement"`
	Name     string     `gorm:"type:text;not null;unique"`
	Students []Students `gorm:"many2many:students_courses;associationForeignKey:StudentID"`
}

type CourseWithStudentCount struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	StudentCount int64  `json:"student_count"`
}
