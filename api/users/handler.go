package users

import (
	"net/http"
	"strconv"
	"log"
	"github.com/gin-gonic/gin"
	util "immy-api/util"
)

type UserHandler struct {
	UserRepo *UserRepo
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, err := h.UserRepo.ListUsers(offset, limit)
	
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
	} else {
		util.OK(c, res)
	}
}

func (h* UserHandler) CreateUser(c *gin.Context) {
	var dto CreateUserDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	res, err := h.UserRepo.CreateUser(dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
	} else {
		util.Created(c, res.ID)
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.NotFound(c, "User", idStr)
	}
	res, err := h.UserRepo.GetUser(id)
	if err != nil {
		util.NotFound(c, "User", id)
	} else {
		util.OK(c, res)
	}
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.NotFound(c, "User", idStr)
	}
	
	var dto UpdateUserDTO
	err = c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}
	
	res, err := h.UserRepo.UpdateUser(id, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
	} else {
		util.OK(c, res)
	}
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.NotFound(c, "User", idStr)
	}
	
	err = h.UserRepo.DeleteUser(id)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
	} else {
		util.OK(c, id)
	}
}