package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseDTO struct {
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

var nilResponse ResponseDTO = ResponseDTO{}

var httpStatusMap = map[string]int{
	"success": http.StatusOK,
	"error":   http.StatusBadRequest,
}

func NilResponse() ResponseDTO {
	return nilResponse
}

func ErrorResponse(err error, errorCode int) (ResponseDTO, error, int) {
	return ResponseDTO{
		Msg: err.Error(),
	}, err, errorCode
}

func SuccessResponse(data interface{}) (ResponseDTO, error, int) {
	return ResponseDTO{
		Data: data,
	}, nil, 200
}

func SendResponse(c *gin.Context, response interface{}) {
	sendHttpResponse(c, response, http.StatusOK)
}

func SendErrorResponse(c *gin.Context, msg string, statusCode int) {
	data := gin.H{
		"msg": msg,
	}
	sendHttpResponse(c, data, statusCode)
}

func sendHttpResponse(c *gin.Context, data interface{}, httpStatus int) {
	c.JSON(httpStatus, data)
}
