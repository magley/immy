package util

import (
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, Response{
		Success: true,
	})
}

func OKPaged(c *gin.Context, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: data,
		Meta: meta,
	})
}

func Created(c *gin.Context, id uint) {
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func Fail(c *gin.Context, status int, code, message string) {
	fmt.Println(status, code, message)
	c.JSON(status, Response{
		Success: false,
		Error:   &ErrorInfo{Code: code, Message: message},
	})
}

func NotFound(c *gin.Context, what string, identifier interface{}) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Error:   &ErrorInfo{
			Code: "NOT_FOUND",
			Message: fmt.Sprintf("Could not find %s by '%v'", what, identifier),
		},
	})
}

func Unauthorized(c *gin.Context, err error) {
	errInfo := &ErrorInfo{
		Code: "UNAUTHORIZED",
	}

	if err != nil {
		errInfo.Message = err.Error()
	}

	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Error:   errInfo,
	})
}