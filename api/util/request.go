package util

import (
	"strconv"
	"strings"

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

func GetJwtToken(c *gin.Context) (string) {
	authVal := c.GetHeader("Authorization")

	if strings.HasPrefix(authVal, "Bearer ") {
		return authVal[7:]
	}

	return ""
}

func GetJwt(c *gin.Context) (*JWTClaims, error) {
	token := GetJwtToken(c)
	if token == "" {
		return nil, nil
	}
	return ValidateJWT(token)
}

func RequireRole(c *gin.Context, role string) (*JWTClaims, bool) {
	jwt, err := GetJwt(c)
	if err != nil {
		Unauthorized(c, err)
		return jwt, false
	}

	err = jwt.RequireRole(role)
	if err != nil {
		Unauthorized(c, err)
		return jwt, false
	}

	return jwt, true
}

func RequireRoleAny(c *gin.Context, roles []string) (*JWTClaims, bool) {
	jwt, err := GetJwt(c)
	if err != nil {
		Unauthorized(c, err)
		return jwt, false
	}

	err = jwt.RequireRoleAny(roles)
	if err != nil {
		Unauthorized(c, err)
		return jwt, false
	}

	return jwt, true
}
