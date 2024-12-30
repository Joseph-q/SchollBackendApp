package courses

import (
	"net/http"

	"github.com/juseph-q/SchoolPr/internal/shared/handler"
)

var CourseNameIsRegister = &CourseErrors{msg: "El nombre del curso ya existe"}
var CourseNotFound = &CourseErrors{msg: "El curso no existe"}

type CourseErrors struct {
	msg string
}

func (e *CourseErrors) Error() string {
	return e.msg
}

func NewCourseErrorNotFoundResponse(err string) *handler.Response {
	return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err, nil)
}

func NewCourseErrorNameRegisterResponse(err string) *handler.Response {
	return handler.NewErrorResponse(http.StatusBadRequest, handler.DuplicateEntry, err, nil)
}
