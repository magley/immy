package handler

import (
	util "immy-api/util"
	"net/http"

	"github.com/gin-gonic/gin"

	"immy-api/model"
	"immy-api/service"
)

type BanAppealHandler struct {
	BanAppealService *service.BanAppealService
	UserService *service.UserService
}

func (h *BanAppealHandler) ListBanAppeals(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin, model.UserRoleModerator})
	if !ok {
		return
	}

	offset, limit := util.GetOffsetLimit(c)
	res, total, err := h.BanAppealService.ListBanAppeals(offset, limit)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OKPaged(c, res, util.MetaPage(limit, offset, total))
		return
	}
}

func (h* BanAppealHandler) CreateBanAppeal(c *gin.Context) {
	var dto model.CreateBanAppealDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}

	res, err := h.BanAppealService.CreateBanAppeal(dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
		return
	}
}

func (h *BanAppealHandler) GetBanAppeal(c *gin.Context) {
	banappealId, ok := util.ParamUintSafe(c, "id", "BanAppeal")
	if !ok {
		return
	}

	res, err := h.BanAppealService.GetBanAppeal(banappealId)
	if err != nil {
		util.NotFound(c, "BanAppeal", banappealId)
		return
	} else {
		util.OK(c, res)
	}
}

func (h *BanAppealHandler) UpdateBanAppeal(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin, model.UserRoleModerator})
	if !ok {
		return
	}

	banappealId, ok := util.ParamUintSafe(c, "id", "BanAppeal")
	if !ok {
		return
	}

	var dto model.UpdateBanAppealDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}

	jwt, err := util.GetJwt(c)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "JWT_FAIL", err.Error())
		return
	}

	user, err := h.UserService.GetUser(jwt.Id)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UNAUTHENTICATED", err.Error())
		return
	}

	res, err := h.BanAppealService.UpdateBanAppeal(banappealId, dto, user)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
	}
}

func (h *BanAppealHandler) DeleteBanAppeal(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin, model.UserRoleModerator})
	if !ok {
		return
	}

	banappealId, ok := util.ParamUintSafe(c, "id", "BanAppeal")
	if !ok {
		return
	}

	err := h.BanAppealService.DeleteBanAppeal(banappealId)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, banappealId)
	}
}

func (h *BanAppealHandler) GetBanAppealsOfBan(c *gin.Context) {
	banID, ok := util.ParamUintSafe(c, "id", "Ban")
	if !ok {
		return
	}

	res, err := h.BanAppealService.GetBanAppealsOfBan(banID)
	if err != nil {
		util.NotFound(c, "BanAppeal", banID)
		return
	} else {
		util.OK(c, res)
	}
}

func (h *BanAppealHandler) CanAppealBan(c *gin.Context) {
	banID, ok := util.ParamUintSafe(c, "id", "Ban")
	if !ok {
		return
	}

	res, err := h.BanAppealService.CanAppealBan(banID)
	if err != nil {
		util.NotFound(c, "BanAppeal", banID)
		return
	} else {
		util.OK(c, res)
	}
}
