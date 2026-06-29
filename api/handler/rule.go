package handler

import (
	util "immy-api/util"
	"net/http"

	"github.com/gin-gonic/gin"

	"immy-api/model"
	"immy-api/service"
)

type RuleHandler struct {
	RuleService *service.RuleService
}

func (h *RuleHandler) ListRules(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, total, err := h.RuleService.ListRules(offset, limit)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OKPaged(c, res, util.MetaPage(limit, offset, total))
		return
	}
}

func (h* RuleHandler) CreateRule(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin})
	if !ok {
		return
	}

	var dto model.CreateRuleDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}

	_, err = util.GetJwt(c)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "JWT_FAIL", err.Error())
		return
	}

	res, err := h.RuleService.CreateRule(dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
		return
	}
}

func (h *RuleHandler) GetRule(c *gin.Context) {
	RuleId, ok := util.ParamUintSafe(c, "id", "Rule")
	if !ok {
		return
	}

	res, err := h.RuleService.GetRule(RuleId)
	if err != nil {
		util.NotFound(c, "Rule", RuleId)
		return
	} else {
		util.OK(c, res)
	}
}

func (h *RuleHandler) UpdateRule(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin})
	if !ok {
		return
	}

	RuleId, ok := util.ParamUintSafe(c, "id", "Rule")
	if !ok {
		return
	}

	var dto model.UpdateRuleDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}

	res, err := h.RuleService.UpdateRule(RuleId, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
	}
}

func (h *RuleHandler) DeleteRule(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin})
	if !ok {
		return
	}

	RuleId, ok := util.ParamUintSafe(c, "id", "Rule")
	if !ok {
		return
	}

	err := h.RuleService.DeleteRule(RuleId)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, RuleId)
	}
}

// ===========================================================================
// RULE <--> BOARD

func (h* RuleHandler) CreateRuleBoard(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin})
	if !ok {
		return
	}

	var dto model.CreateRuleBoardDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}

	_, err = util.GetJwt(c)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "JWT_FAIL", err.Error())
		return
	}

	res, err := h.RuleService.CreateRuleBoard(dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.BoardID) // TODO: This is not the primary key of a rule-board.
		return
	}
}

func (h *RuleHandler) DeleteRuleBoard(c *gin.Context) {
	_, ok := util.RequireRoleAny(c, []string{model.UserRoleAdmin})
	if !ok {
		return
	}

	boardId, ok := util.ParamUintSafe(c, "boardId", "Rule")
	if !ok {
		return
	}

	ruleId, ok := util.ParamUintSafe(c, "ruleId", "Rule")
	if !ok {
		return
	}

	err := h.RuleService.DeleteRuleBoard(boardId, ruleId)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, nil)
	}
}

func (h *RuleHandler) ListAllRuleBoards(c *gin.Context) {
	res, total, err := h.RuleService.ListAllRuleBoards()

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OKPaged(c, res, util.MetaPage(0, 0, total))
		return
	}
}

func (h *RuleHandler) ListAllRulesOfBoard(c *gin.Context) {
	boardId, ok := util.ParamUintSafe(c, "boardId", "Rule")
	if !ok {
		return
	}

	res, total, err := h.RuleService.ListAllRulesOfBoard(boardId)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OKPaged(c, res, util.MetaPage(0, 0, total))
		return
	}
}

func (h *RuleHandler) ListAllBoardsOfRule(c *gin.Context) {
	ruleId, ok := util.ParamUintSafe(c, "ruleId", "Rule")
	if !ok {
		return
	}

	res, total, err := h.RuleService.ListAllBoardsOfRule(ruleId)

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OKPaged(c, res, util.MetaPage(0, 0, total))
		return
	}
}