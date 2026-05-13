package handler


import (
	"net/http"
	"github.com/gin-gonic/gin"
	util "immy-api/util"
		
	"immy-api/service"
	"immy-api/model"
)

type UserHandler struct {
	UserService *service.UserService
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, err := h.UserService.ListUsers(offset, limit)
	
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h* UserHandler) CreateUser(c *gin.Context) {
	var dto model.CreateUserDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	res, err := h.UserService.CreateUser(dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
		return
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userId, ok := util.ParamUintSafe(c, "id", "User")
	if !ok {
		return
	}

	res, err := h.UserService.GetUser(userId)
	if err != nil {
		util.NotFound(c, "User", userId)
		return
	} else {
		util.OK(c, res)
	}
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userId, ok := util.ParamUintSafe(c, "id", "User")
	if !ok {
		return
	}

	var dto model.UpdateUserDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}
	
	res, err := h.UserService.UpdateUser(userId, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
	}
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userId, ok := util.ParamUintSafe(c, "id", "User")
	if !ok {
		return
	}

	err := h.UserService.DeleteUser(userId)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, userId)
	}
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var dto model.LoginUserDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	result, err := h.UserService.LoginUser(dto)
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, "LOGIN_FAIL", err.Error())
		return
	} else {
		util.OK(c, result)
		return
	}
}