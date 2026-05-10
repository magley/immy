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