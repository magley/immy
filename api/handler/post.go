package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	util "immy-api/util"
	
	"immy-api/service"
	"immy-api/model"
)

type PostHandler struct {
	PostService *service.PostService
}


func (h *PostHandler) ListPosts(c *gin.Context) {
	offset, limit := util.GetOffsetLimit(c)
	res, err := h.PostService.ListPosts(offset, limit)
	
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var dto model.CreatePostDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "BAD_JSON", err.Error())
		return
	}
	
	res, err := h.PostService.CreatePost(dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "CREATE_FAIL", err.Error())
		return
	} else {
		util.Created(c, res.ID)
		return
	}
}

func (h *PostHandler) GetPostByNum(c *gin.Context) {
	boardCode := c.Param("boardCode")
	postNum, ok := util.ParamUintSafe(c, "postNum", "Post")
	if !ok {
		return
	} 
	
	res, err := h.PostService.GetPostByNum(boardCode, postNum)
	if err != nil {
		util.NotFound(c, "Post", postNum)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *PostHandler) GetPostsByThread(c *gin.Context) {
	threadId, ok := util.ParamUintSafe(c, "threadId", "Thread")
	if !ok {
		return
	} 
	
	res, err := h.PostService.GetPostsByThread(threadId)
	if err != nil {
		util.NotFound(c, "Posts of thread", threadId)
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	postId, ok := util.ParamUintSafe(c, "id", "Post")
	if !ok {
		return
	}
	
	var dto model.UpdatePostDTO
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "ERROR", err.Error())
		return
	}
	
	res, err := h.PostService.UpdatePost(postId, dto)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "UPDATE_FAILED", err.Error())
		return
	} else {
		util.OK(c, res)
		return
	}
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	postId, ok := util.ParamUintSafe(c, "id", "Post")
	if !ok {
		return
	}
			
	err := h.PostService.DeletePost(postId)
	if err != nil {
		util.Fail(c, http.StatusBadRequest, "DELETE_FAIL", err.Error())
		return
	} else {
		util.OK(c, postId)
		return
	}
}