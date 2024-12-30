package assistanceController

import (
	"errors"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juseph-q/SchoolPr/internal/assistance"
	assistanceServices "github.com/juseph-q/SchoolPr/internal/assistance/services"
	"github.com/juseph-q/SchoolPr/internal/courses"
	coursesService "github.com/juseph-q/SchoolPr/internal/courses/service"
	"github.com/juseph-q/SchoolPr/internal/shared/handler"
	"github.com/juseph-q/SchoolPr/internal/shared/utils"
)

type AssistanceController struct {
	assitanceService *assistanceServices.AssistanceService
	courseService    *coursesService.CourseService
}

func NewAssistanceController(assistance *assistanceServices.AssistanceService, course *coursesService.CourseService) *AssistanceController {
	return &AssistanceController{
		assitanceService: assistance,
		courseService:    course,
	}
}

func (h *AssistanceController) getAssitanceByCourseId(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		id, err := utils.ParseIdParam(c)
		if err != nil {
			handler.NewInvalidParamResponse("Invalid id")
		}
		//Comprobar si existe el curso id
		course, err := h.courseService.FindCourseById(uint(id), c.Request.Context())

		if err != nil {
			if errors.Is(err, courses.CourseNotFound) {
				return courses.NewCourseErrorNotFoundResponse(err.Error())
			}
			return handler.NewInternalErrorResponse(err)
		}

		var queryReq assistance.QueryParamsByCourseId
		if err := c.ShouldBindQuery(&queryReq); err != nil {
			return handler.NewInvalidParamResponse(err.Error())
		}
		if queryReq.Limit <= 0 {
			queryReq.Limit = 100
		}
		if queryReq.Page <= 0 {
			queryReq.Page = 1
		}

		studentsAssisted, total, err := h.assitanceService.StudentsAssitanceByCourseId(course, &queryReq, c.Request.Context())

		if err != nil {
			if errors.Is(err, assistance.NotFountRegister) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, assistance.NewCourseResponseAssistance(studentsAssisted, queryReq.Date, handler.MetadataPage{
			Page:      queryReq.Page,
			PageSize:  queryReq.Limit,
			PageCount: int(math.Ceil(float64(total) / float64(queryReq.Limit))),
			Total:     total,
		}))

	})
}

func (h *AssistanceController) getHistorialAssitance(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		var query assistance.QueryParamsHistorialAssitance
		if err := c.ShouldBindQuery(&query); err != nil {
			return handler.NewInvalidParamResponse("Invalid query params")
		}
		response, err := h.assitanceService.HistorialAssistances(c.Request.Context(), &query)
		if err != nil {
			if !errors.Is(err, handler.ServerError) {
				if errors.Is(err, courses.CourseNotFound) {
					return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
				}
				if errors.Is(err, assistance.NotFountRegister) {
					return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
				}
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, assistance.NewResponseHistorial(response))
	})
}

func (h *AssistanceController) getAssitancesStudent(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		id, err := utils.ParseIdParam(c)
		if err != nil {
			handler.NewInvalidParamResponse("Invalid id")
		}

		var queryReq assistance.QueryParamsAssistanceStudentById
		if err := c.ShouldBindQuery(&queryReq); err != nil {
			return handler.NewInvalidParamResponse(err.Error())
		}

		assistanceDb, err := h.assitanceService.FindStudentAssistance(c.Request.Context(), id, &queryReq)
		if err != nil {
			if errors.Is(err, assistance.NotFountRegister) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, assistance.NewResponseAssistancesStudent(uint(id), assistanceDb))

	})
}
func (h *AssistanceController) getStudentCourseAssisted(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		id, err := utils.ParseIdParam(c)
		if err != nil {
			handler.NewInvalidParamResponse("Invalid id")
		}
		assistanceDb, err := h.assitanceService.FindStudentCoursesAssisted(c.Request.Context(), id)

		if err != nil {
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, assistanceDb)

	})
}
func (h *AssistanceController) getHistorialCourse(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		var queryReq assistance.QueryParamsHistorialAssitanceSummary
		if err := c.ShouldBindQuery(&queryReq); err != nil {
			return handler.NewInvalidParamResponse("Invalid query params")

		}

		response, err := h.assitanceService.HistorialAssistancesSumary(c.Request.Context(), &queryReq)

		if err != nil {
			if errors.Is(err, assistance.NotFountRegister) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, assistance.NewResponseAssistancesSumary(response, handler.MetadataPage{
			Page:      1,
			PageSize:  1,
			PageCount: 1,
			Total:     1,
		}))
	})
}

func HandleRoutes(r *gin.Engine, s *AssistanceController) {
	r.GET("/assistances/course/:id", s.getAssitanceByCourseId)

	r.GET("/assistances/students/:id", s.getAssitancesStudent)

	r.GET("/assistances/historial", s.getHistorialAssitance)
	r.GET("/assistances/historial/all", s.getHistorialCourse)

	r.GET("/assistances/studentsAssistedCourses/:id", s.getStudentCourseAssisted)

}
