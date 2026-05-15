package datatransfers

import "github.com/gin-gonic/gin"

// here we will create a standarization format for respnse API

// successs response
type SuccessResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// standard error response
type ErrorResponse struct {
	StatusCode int      `json:"statusCode"`
	Message    string   `json:"message"`
	Error      []string `json:"error"`
}

// success helper
func SuccessRes(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, SuccessResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func ErrorRes(c *gin.Context, statusCode int, message string, errors ...string) {
	c.JSON(statusCode, ErrorResponse{
		StatusCode: statusCode,
		Message:    message,
		Error:      errors,
	})
}
