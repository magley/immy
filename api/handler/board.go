package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	util "immy-api/util"
	
	"immy-api/service"
	"immy-api/model"
)

type BoardHandler struct {
	BoardService *service.BoardService
}


func (h *BoardHandler) ListBoards(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, err := h.BoardService.ListBoards(offset, limit)
	
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	var dto model.CreateBoardDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	res, err := h.BoardService.CreateBoard(dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
		return
	}
}

func (h *BoardHandler) GetBoard(c *gin.Context) {
	boardCode := c.Param("code")
	
	res, err := h.BoardService.GetBoardByCode(boardCode)
	if err != nil {
		util.NotFound(c, "Board", boardCode)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	boardCode := c.Param("code")

	var dto model.UpdateBoardDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}
	
	res, err := h.BoardService.UpdateBoard(boardCode, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	boardCode := c.Param("code")
		
	err := h.BoardService.DeleteBoard(boardCode)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, boardCode)
		return
	}
}