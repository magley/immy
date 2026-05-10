package boards

import (
	"net/http"
	"github.com/gin-gonic/gin"
	util "immy-api/util"
)

type BoardHandler struct {
	BoardRepo *BoardRepo
}


func (h *BoardHandler) ListBoards(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, err := h.BoardRepo.ListBoards(offset, limit)
	
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
	} else {
		util.OK(c, res)
	}
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	var dto CreateBoardDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	res, err := h.BoardRepo.CreateBoard(dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
	} else {
		util.Created(c, res.ID)
	}
}

func (h *BoardHandler) GetBoard(c *gin.Context) {
	boardCode := c.Param("code")
	
	res, err := h.BoardRepo.GetBoard(boardCode)
	if err != nil {
		util.NotFound(c, "Board", boardCode)
	} else {
		util.OK(c, res)
	}
}

func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	boardCode := c.Param("code")

	var dto UpdateBoardDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}
	
	res, err := h.BoardRepo.UpdateBoard(boardCode, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
	} else {
		util.OK(c, res)
	}
}

func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	boardCode := c.Param("code")
	
	err := h.BoardRepo.DeleteBoard(boardCode)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
	} else {
		util.OK(c, boardCode)
	}
}