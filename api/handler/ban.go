package handler


import (
	"net/http"
	"github.com/gin-gonic/gin"
	util "immy-api/util"

	"immy-api/service"
	"immy-api/model"
)

type BanHandler struct {
	BanService *service.BanService
	UserService *service.UserService
}

func (h *BanHandler) ListBans(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, err := h.BanService.ListBans(offset, limit)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *BanHandler) ListBansForAdmin(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin, model.UserRoleModerator})
	if !ok {
		return
	}

	offset, limit := util.GetOffsetLimit(c)
	res, err := h.BanService.ListBansForAdmin(offset, limit)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}


func (h *BanHandler) GetMyBans(c *gin.Context) {
	ip := c.Copy().ClientIP()
	res, err := h.BanService.GetBansOfIp(ip)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "GET_BANS_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h* BanHandler) CreateBan(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin, model.UserRoleModerator})
	if !ok {
		return
	}

	var dto model.CreateBanDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
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

	res, err := h.BanService.CreateBan(dto, user)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
		return
	}
}

func (h *BanHandler) GetBan(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin, model.UserRoleModerator})
	if !ok {
		return
	}

	banId, ok := util.ParamUintSafe(c, "id", "Ban")
	if !ok {
		return
	}

	res, err := h.BanService.GetBan(banId)
	if err != nil {
		util.NotFound(c, "Ban", banId)
		return
	} else {
		util.OK(c, res)
	}
}

func (h *BanHandler) GetBanForAdmin(c *gin.Context) {
	banId, ok := util.ParamUintSafe(c, "id", "Ban")
	if !ok {
		return
	}

	res, err := h.BanService.GetBanForAdmin(banId)
	if err != nil {
		util.NotFound(c, "Ban", banId)
		return
	} else {
		util.OK(c, res)
	}
}

func (h *BanHandler) UpdateBan(c *gin.Context) {
	// No authorization because as of right now, only the "seen" field can be
	// updated. Should the admins/mods in the future be able to update other
	// stuff, split "mark seen" into its own thing, or do a more clever
	// access control.

	banId, ok := util.ParamUintSafe(c, "id", "Ban")
	if !ok {
		return
	}

	var dto model.UpdateBanDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}

	res, err := h.BanService.UpdateBan(banId, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
	}
}

func (h *BanHandler) DeleteBan(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin})
	if !ok {
		return
	}

	banId, ok := util.ParamUintSafe(c, "id", "Ban")
	if !ok {
		return
	}

	err := h.BanService.DeleteBan(banId)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, banId)
	}
}