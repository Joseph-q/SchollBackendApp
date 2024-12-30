package handler

import "net/http"

type Response struct {
	StatusCode int
	Data       interface{}
	Err        error
}

type MetadataPage struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"pageSize"`
	PageCount int   `json:"pageCount"`
	Total     int64 `json:"total"`
}

func NewSuccessResponse(statusCode int, data interface{}) *Response {
	return &Response{
		StatusCode: statusCode,
		Data:       data,
	}
}

func NewErrorResponse(statusCode int, code ErrorCode, message string, details interface{}) *Response {
	return &Response{
		StatusCode: statusCode,
		Err: &ErrorResponse{
			Code:    code,
			Message: message,
			Errors:  details,
		},
	}
}

func NewInternalErrorResponse(err error) *Response {
	return &Response{
		Err: err,
	}
}

func NewInvalidQueryResponse(message string) *Response {
	return NewErrorResponse(http.StatusBadRequest, InvalidUriValue, message, nil)
}

func NewInvalidParamResponse(message string) *Response {
	return NewErrorResponse(http.StatusBadRequest, InvalidParamValue, message, nil)
}
