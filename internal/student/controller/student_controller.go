package studentController

import (
	"errors"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juseph-q/SchoolPr/internal/config"

	"github.com/juseph-q/SchoolPr/internal/shared/handler"
	"github.com/juseph-q/SchoolPr/internal/shared/middleware"
	"github.com/juseph-q/SchoolPr/internal/shared/utils"
	"github.com/juseph-q/SchoolPr/internal/student"
	studentService "github.com/juseph-q/SchoolPr/internal/student/services"
)

type StudentController struct {
	service *studentService.StudentService
}

func NewStudentController(s *studentService.StudentService) *StudentController {
	return &StudentController{
		service: s,
	}
}

func (s *StudentController) getStudent(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		id, err := utils.ParseIdParam(c)
		if err != nil {
			return handler.NewInvalidParamResponse("Invalid id")
		}
		var queryReq student.QueryGetStudentById

		if err = c.ShouldBindQuery(&queryReq); err != nil {
			return handler.NewInvalidQueryResponse("Invalid query")
		}

		studentDB, err := s.service.FindStudentById(id, &queryReq, c.Request.Context())

		if err != nil {
			if !errors.Is(err, handler.ServerError) {
				if errors.Is(err, student.UserNotFound) {
					return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, student.UserNotFound.Error(), nil)

				}
			}

			return handler.NewErrorResponse(http.StatusInternalServerError, handler.InternalServerError, handler.ServerError.Error(), nil)
		}

		return handler.NewSuccessResponse(http.StatusOK, student.NewStudentResponse(studentDB))
	})
}

func (s *StudentController) getStudents(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		var queryReq student.QueryGetStudents
		if err := c.ShouldBindQuery(&queryReq); err != nil {
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidQueryValue, "invalid query params", nil)
		}
		if queryReq.Limit <= 0 {
			queryReq.Limit = 50
		}
		if queryReq.Page <= 0 {
			queryReq.Page = 1
		}

		studentsDb, total, err := s.service.FindStudents(queryReq, c.Request.Context())

		if studentsDb == nil && err != nil {
			return handler.NewInternalErrorResponse(handler.ServerError)
		}

		return handler.NewSuccessResponse(http.StatusOK, student.NewStudentsResponse(studentsDb, handler.MetadataPage{
			Page:      queryReq.Page,
			PageSize:  queryReq.Limit,
			PageCount: int(math.Ceil(float64(total) / float64(queryReq.Limit))),
			Total:     total,
		}, queryReq))

	})
}

func (s *StudentController) createStudent(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {

		var studentTC student.CreateStudent

		if err := c.ShouldBindJSON(&studentTC); err != nil {
			return handler.NewInvalidParamResponse(err.Error())
		}

		if studentTC.Email == nil && studentTC.Number == nil {
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidBodyValue, "Email o numero deberian tener valor", nil)
		}

		_, err := s.service.CreateStudent(studentTC, c.Request.Context())

		if err != nil {
			if errors.Is(err, student.CourseNotValid) {
				return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidBodyValue, err.Error(), nil)
			}
			if errors.Is(err, student.EmailOrNumberAlreadyRegister) {
				return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidBodyValue, err.Error(), nil)
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusCreated, nil)
	})
}

func (s *StudentController) updateStudent(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		id, err := utils.ParseIdParam(c)
		if err != nil {
			return handler.NewInvalidParamResponse("Invalid id")
		}
		var studentTU student.UpdateStudent
		if err := c.ShouldBindBodyWithJSON(&studentTU); err != nil {
			return handler.NewInvalidParamResponse(err.Error())
		}

		studentDb, err := s.service.FindStudentById(id, nil, c.Request.Context())
		if err != nil {
			if errors.Is(err, student.UserNotFound) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		err = s.service.UpdateStudent(student.NewStudentUpdateDb(&studentTU, studentDb), c.Request.Context())

		if err != nil {
			if errors.Is(err, student.UserNotFound) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
			}
			if errors.Is(err, student.EmailOrNumberAlreadyRegister) {
				return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidBodyValue, err.Error(), nil)
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusNoContent, nil)
	})
}

func (s *StudentController) deleteStudent(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		id, err := utils.ParseIdParam(c)
		if err != nil {
			return handler.NewInvalidParamResponse("Invalid id")
		}
		studentDb, err := s.service.FindStudentById(id, nil, c)

		if err != nil {
			if errors.Is(err, student.UserNotFound) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		err = s.service.DeleteStudent(studentDb, c.Request.Context())
		if err != nil {
			if errors.Is(err, student.UserNotFound) {
				return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
			}
			return handler.NewInternalErrorResponse(err)
		}
		//borar indice

		return handler.NewSuccessResponse(http.StatusNoContent, nil)
	})
}

func (s *StudentController) searchStudent(c *gin.Context) {
	handler.HandleRequest(c, func(c *gin.Context) *handler.Response {
		var query student.QuerySearchStudent
		if err := c.ShouldBindQuery(&query); err != nil {
			return handler.NewErrorResponse(http.StatusBadRequest, handler.InvalidQueryValue, "invalid query params", nil)
		}

		studentSearch, totalHits, err := s.service.SearchStudent(&query, c.Request.Context())

		if err != nil {
			if !errors.Is(err, handler.ServerError) {
				if errors.Is(err, student.UserNotFound) {
					return handler.NewErrorResponse(http.StatusNotFound, handler.NotFoundEntity, err.Error(), nil)
				}
			}
			return handler.NewInternalErrorResponse(err)
		}

		return handler.NewSuccessResponse(http.StatusOK, student.NewStudentSearchReponse(studentSearch, totalHits, query.Limit, query.Query))
	})
}

func HandleRoutes(router *gin.Engine, c *StudentController, cf *config.Config) {

	router.Use(middleware.TimeOutMiddleware(cf.Server.MaxWriteTimeout))

	router.GET("/students/:id", c.getStudent)
	router.GET("/students", c.getStudents)

	router.POST("/students", c.createStudent)
	router.PUT("/students/:id", c.updateStudent)

	router.DELETE("/students/:id", c.deleteStudent)

	// search
	router.GET("/students/search", c.searchStudent)

}
