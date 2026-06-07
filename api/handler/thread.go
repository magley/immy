package handler

import (
	"immy-api/util"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"

	"immy-api/model"
	"immy-api/service"
)

type ThreadHandler struct {
	ThreadService 	*service.ThreadService
	BoardService 	*service.BoardService
	PostService 	*service.PostService
}

func (h *ThreadHandler) ListThreads(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, err := h.ThreadService.ListThreads(offset, limit)
	
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
	
	res, err := h.ThreadService.ListThreadsOfBoard(boardCode, offset, limit)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) CreateThread(c *gin.Context) {
	var dto model.CreateThreadDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	res, err := h.ThreadService.CreateThread(dto, c.ClientIP())
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
		return
	}
}

func (h *ThreadHandler) GetThread(c *gin.Context) {
	threadId, ok := util.ParamUintSafe(c, "id", "Thread")
	if !ok {
		return
	}
	
	res, err := h.ThreadService.GetThread(threadId)
	if err != nil {
		util.NotFound(c, "Thread", threadId)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) GetFullThread(c *gin.Context) {
	threadId, ok := util.ParamUintSafe(c, "id", "Thread")
	if !ok {
		return
	}
	
	res, err := h.ThreadService.GetFullThread(threadId)
	if err != nil {
		util.NotFound(c, "Thread", threadId)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) GetThreadByNum(c *gin.Context) {
	boardCode := c.Param("boardCode")
	threadNum, ok := util.ParamUintSafe(c, "num", "Thread")
	if !ok {
		return
	}
	
	res, err := h.ThreadService.GetThreadByNum(boardCode, threadNum)
	if err != nil {
		util.NotFound(c, "Thread", threadNum)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) GetFullThreadByNum(c *gin.Context) {
	boardCode := c.Param("boardCode")
	threadNum, ok := util.ParamUintSafe(c, "num", "Thread")
	if !ok {
		return
	}
	
	res, err := h.ThreadService.GetFullThreadByNum(boardCode, threadNum)
	if err != nil {
		util.NotFound(c, "Thread", threadNum)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) GetThreadsForCatalog(c *gin.Context) {
	boardCode := c.Param("boardCode")

	res, err := h.ThreadService.GetThreadsForCatalog(boardCode)
	if err != nil {
		util.NotFound(c, "Threads of board", boardCode)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) GetThreadsForArchive(c *gin.Context) {
	boardCode := c.Param("boardCode")

	res, err := h.ThreadService.GetThreadsForArchive(boardCode)
	if err != nil {
		util.NotFound(c, "Threads of board", boardCode)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) GetThreadsForHome(c *gin.Context) {
	boardCode := c.Param("boardCode")
	page, ok := util.QueryIntSafe(c, "page", "Threads")
	if !ok {
		return
	}
	perPage, ok := util.QueryIntSafe(c, "perPage", "Threads")
	if !ok {
		return
	}
	perThread := 5

	res, totalThreads, err := h.ThreadService.GetThreadsForHome(boardCode, perThread, perPage, page)

	if err != nil {
		util.NotFound(c, "Threads of board", boardCode)
		return
	} else {
		util.OKPaged(c, res, &util.Meta{
			Page: page,
			PerPage: perPage,
			Total: int(totalThreads),
			TotalPages: int(math.Ceil(float64(totalThreads) / float64(perPage))),
		})
		return
	}
}

func (h *ThreadHandler) UpdateThread(c *gin.Context) {
	threadId, ok := util.ParamUintSafe(c, "id", "Thread")
	if !ok {
		return
	}
	
	var dto model.UpdateThreadDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}
	
	res, err := h.ThreadService.UpdateThread(threadId, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *ThreadHandler) DeleteThread(c *gin.Context) {
	threadId, ok := util.ParamUintSafe(c, "id", "Thread")
	if !ok {
		return
	}
	
	err := h.ThreadService.DeleteThread(threadId)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, threadId)
		return
	}
}

func (h *ThreadHandler) ArchiveThread(c *gin.Context) {
	threadId, ok := util.ParamUintSafe(c, "id", "Thread")
	if !ok {
		return
	}

	res, err := h.ThreadService.ArchiveThread(threadId)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

