package ginutils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerWithErrFunc = func(*gin.Context) error
type HandlerFunc = HandlerWithErrFunc

func HandlerWithErr(h HandlerWithErrFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := h(ctx)
		if err != nil {
			var ginErr ErrorResponse
			if !errors.As(err, &ginErr) {
				ginErr = ServerError(err.Error())
			}

			ctx.AbortWithStatusJSON(ginErr.Status, ginErr)
		}
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			errResponse := recoveredErrToErrorResponse(err)
			errResponse.AbortJSON(c)
		}()

		c.Next()
	}
}

func recoveredErrToErrorResponse(e any) ErrorResponse {
	switch err := e.(type) {
	case ErrorResponse:
		return err
	case error:
		return ServerError(err.Error())
	case string:
		return ServerError(err)
	default:
		return ServerError("internal server error")
	}
}

func AbortWithError(c *gin.Context, err ErrorResponse) {
	c.AbortWithStatusJSON(err.Status, err)
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (g ErrorResponse) Error() string {
	return g.Message
}

func (g ErrorResponse) AbortJSON(c *gin.Context) {
	AbortWithError(c, g)
}

func ServerError(message string) ErrorResponse {
	return ErrorResponse{
		Status:  500,
		Message: message,
	}
}

func BadRequest(message string) ErrorResponse {
	return ErrorResponse{
		Status:  400,
		Message: message,
	}
}

func Unauthorized(message string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func Forbidden(message string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusForbidden,
		Message: message,
	}
}

func NotFound(message string) ErrorResponse {
	return ErrorResponse{
		Status:  http.StatusNotFound,
		Message: message,
	}
}
