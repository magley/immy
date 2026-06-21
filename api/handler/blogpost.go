package handler

import (
	util "immy-api/util"
	"net/http"

	"github.com/gin-gonic/gin"

	"immy-api/model"
	"immy-api/service"
)

type BlogpostHandler struct {
	BlogpostService *service.BlogpostService
	UserService *service.UserService
}

func (h *BlogpostHandler) ListBlogposts(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, total, err := h.BlogpostService.ListBlogposts(offset, limit)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OKPaged(c, res, util.MetaPage(limit, offset, total))
		return
	}
}

func (h *BlogpostHandler) ListBlogpostsShort(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, total, err := h.BlogpostService.ListBlogposts(offset, limit)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OKPaged(c, h.BlogpostService.ToShortArr(res), util.MetaPage(limit, offset, total))
		return
	}
}

func (h* BlogpostHandler) CreateBlogpost(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin})
	if !ok {
		return
	}

	var dto model.CreateBlogpostDTO
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

	res, err := h.BlogpostService.CreateBlogpost(dto, user)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
		return
	}
}

func (h *BlogpostHandler) GetBlogpost(c *gin.Context) {
	blogpostId, ok := util.ParamUintSafe(c, "id", "Blogpost")
	if !ok {
		return
	}

	res, err := h.BlogpostService.GetBlogpost(blogpostId)
	if err != nil {
		util.NotFound(c, "Blogpost", blogpostId)
		return
	} else {
		util.OK(c, res)
	}
}

func (h *BlogpostHandler) UpdateBlogpost(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin})
	if !ok {
		return
	}

	blogpostId, ok := util.ParamUintSafe(c, "id", "Blogpost")
	if !ok {
		return
	}

	var dto model.UpdateBlogpostDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}

	res, err := h.BlogpostService.UpdateBlogpost(blogpostId, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
	}
}

func (h *BlogpostHandler) DeleteBlogpost(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin})
	if !ok {
		return
	}

	blogpostId, ok := util.ParamUintSafe(c, "id", "Blogpost")
	if !ok {
		return
	}

	err := h.BlogpostService.DeleteBlogpost(blogpostId)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, blogpostId)
	}
}