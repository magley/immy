package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)


func GetOffsetLimit(c *gin.Context) (offset int, limit int) {
	offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit > 100 { limit = 100 }
	
	return offset, limit
}

func ParamUint(c *gin.Context, paramName string) (uint, error) {
	valueStr := c.Param(paramName)
	value, err := strconv.ParseUint(valueStr, 10, 32)
	return uint(value), err
}

func ParamUintSafe(c *gin.Context, paramName string, what string) (uint, bool) {
	valueStr := c.Param(paramName)
	value, err := strconv.ParseUint(valueStr, 10, 32)
	if err != nil {
		NotFound(c, what, valueStr)
		return 0, false
	}
	
	return uint(value), true
}

func ParamInt(c *gin.Context, paramName string) (int, error) {
	valueStr := c.Param(paramName)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	return int(value), err
}

func ParamIntSafe(c *gin.Context, paramName string, what string) (int, bool) {
	valueStr := c.Param(paramName)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		NotFound(c, what, valueStr)
		return 0, false
	}

	return int(value), true
}


func QueryInt(c *gin.Context, paramName string) (int, error) {
	valueStr := c.Query(paramName)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	return int(value), err
}

func QueryIntSafe(c *gin.Context, paramName string, what string) (int, bool) {
	valueStr := c.Query(paramName)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		NotFound(c, what, valueStr)
		return 0, false
	}

	return int(value), true
}