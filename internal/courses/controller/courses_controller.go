package coursesController

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/juseph-q/SchoolPr/internal/courses"
	coursesService "github.com/juseph-q/SchoolPr/internal/courses/service"
	"github.com/juseph-q/SchoolPr/internal/shared/handler"
	"github.com/juseph-q/SchoolPr/internal/shared/utils"
)

type CourseController struct {
	service *coursesService.CourseService
}

func NewCourseController(service *coursesService.CourseService) *CourseController {
	return &CourseController{
		service: service,
	}
}

func (s *CourseController) getCourses(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {

		var query courses.QueryGetCourses
		if err := c.ShouldBindQuery(&query); err != nil {
			return handler.NewInvalidParamResponse("Invalid query params")
		}

		if query.Limit <= 0 {
			query.Limit = 20
		}

		if query.Page <= 0 {
			query.Page = 1

		}

		fmt.Println(query.Page)

		courseDb, total, err := s.service.FindCourses(&query, c.Request.Context())

		if err != nil {
			return handler.NewInternalErrorResponse(handler.ServerError)
		}

		return handler.NewSuccessResponse(http.StatusOK, courses.NewCoursesResponse(courseDb, handler.MetadataPage{
			Page:      query.Page,
			PageSize:  query.Limit,
			PageCount: int(math.Ceil((float64(total)) / float64(query.Limit))),
			Total:     total,
		}))
	})

}

func (s *CourseController) createCourse(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		var courseTC courses.CreateOrUpdateCourse
		if err := c.ShouldBindBodyWithJSON(&courseTC); err != nil {
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidBodyValue, err.Error(), nil)
		}

		courseDb, err := s.service.CreateCourse(&courseTC, c.Request.Context())

		if err != nil {
			if !errors.Is(err, handler.ServerError) {
				if errors.Is(err, courses.CourseNameIsRegister) {
					return courses.NewCourseErrorNameRegisterResponse(err.Error())
				}
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusCreated, courses.NewCourseResponse(courseDb))

	})
}

func (s *CourseController) updateCourse(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		id, err := utils.ParseIdParam(c)
		if err != nil {
			handler.NewInvalidParamResponse("Invalid id")
		}
		courseDb, err := s.service.FindCourseById(uint(id), c.Request.Context())

		if err != nil {
			if errors.Is(err, courses.CourseNotFound) {
				courses.NewCourseErrorNotFoundResponse(err.Error())
			}
			return handler.NewInternalErrorResponse(err)
		}

		var courseTU courses.CreateOrUpdateCourse
		if err := c.ShouldBindBodyWithJSON(&courseTU); err != nil {
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidBodyValue, err.Error(), nil)
		}

		if courseTU.Name != "" {
			courseDb.Name = courseTU.Name
		}

		err = s.service.UpdateCourse(c.Request.Context(), courseDb)

		if err != nil {
			if !errors.Is(err, handler.ServerError) {
				if errors.Is(err, courses.CourseNotFound) {
					return courses.NewCourseErrorNotFoundResponse(err.Error())
				}
				if errors.Is(err, courses.CourseNameIsRegister) {
					return courses.NewCourseErrorNameRegisterResponse(err.Error())

				}
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusNoContent, nil)
	})
}

func HandleRoutes(r *gin.Engine, s *CourseController) {
	r.GET("/courses", s.getCourses)

	r.POST("/courses", s.createCourse)
	r.PUT("/courses/:id", s.updateCourse)

}
