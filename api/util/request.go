package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
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