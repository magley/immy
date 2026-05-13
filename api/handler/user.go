package handler


import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	util "immy-api/util"
		
	_ "immy-api/service"
	"immy-api/repo"
	"immy-api/model"
)

type UserHandler struct {
	UserRepo *repo.UserRepo
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, err := h.UserRepo.ListUsers(offset, limit)
	
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
	}
}

func (h* UserHandler) CreateUser(c *gin.Context) {
	var dto model.CreateUserDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	res, err := h.UserRepo.CreateUser(dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.NotFound(c, "User", idStr)
		return
	}
	res, err := h.UserRepo.GetUser(id)
	if err != nil {
		util.NotFound(c, "User", id)
		return
	} else {
		util.OK(c, res)
	}
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.NotFound(c, "User", idStr)
		return
	}
	
	var dto model.UpdateUserDTO
	err = c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}
	
	res, err := h.UserRepo.UpdateUser(id, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
	}
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		util.NotFound(c, "User", idStr)
		return
	}
	
	err = h.UserRepo.DeleteUser(id)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, id)
	}
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var dto model.LoginUserDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	user, err := h.UserRepo.GetUserByName(dto.Username)
	
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, "LOGIN_FAIL", err.Error())
		return
	}
	
	ok := util.CheckPasswordHash(dto.Password, user.Password)
	if !ok {
		util.Fail(c, http.StatusUnauthorized, "LOGIN_FAIL", "Unauthorized")
		return
	}
	
	jwt, err := util.CreateJWT(user.ID, user.Username, string(user.Type))
	if err != nil {
		util.Fail(c, http.StatusUnauthorized, "LOGIN_FAIL", err.Error())
		return
	}
	
	res := model.LoginResponseDTO{
		ID: user.ID,
		Username: user.Username,
		Type: user.Type,
		JWT: jwt,
	}
	util.OK(c, res)
}