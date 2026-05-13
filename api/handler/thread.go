package handler


import (
	"net/http"
	"github.com/gin-gonic/gin"
	util "immy-api/util"
			
	_ "immy-api/service"
	"immy-api/repo"
	"immy-api/model"
)

type ThreadHandler struct {
	ThreadRepo 	*repo.ThreadRepo
	BoardRepo 	*repo.BoardRepo
	PostRepo 	*repo.PostRepo
}

func (h *ThreadHandler) ListThreads(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, err := h.ThreadRepo.ListThreads(offset, limit)
	
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) ListThreadsOfBoard(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	boardCode := c.Param("boardCode")
	
	board, err := h.BoardRepo.GetBoard(boardCode)
	if err != nil {
		util.NotFound(c, "Board", boardCode)
		return
	}
	
	res, err := h.ThreadRepo.ListThreadsOfBoard(board.ID, offset, limit)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}


func (h *ThreadHandler) CreateThread(c *gin.Context) {
	// Get board
	
	
	// Get post number
	
	// Create post
	
	// Create thread
	
	
	var dto model.CreateThreadDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	res, err := h.ThreadRepo.CreateThread(dto, 0)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
		return
	}
}

func (h *ThreadHandler) GetThread(c *gin.Context) {
	threadCode := c.Param("code")
	
	res, err := h.ThreadRepo.GetThread(threadCode)
	if err != nil {
		util.NotFound(c, "Thread", threadCode)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) UpdateThread(c *gin.Context) {
	threadCode := c.Param("code")

	var dto model.UpdateThreadDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}
	
	res, err := h.ThreadRepo.UpdateThread(threadCode, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) DeleteThread(c *gin.Context) {
	threadCode := c.Param("code")
	
	err := h.ThreadRepo.DeleteThread(threadCode)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, threadCode)
		return
	}
}