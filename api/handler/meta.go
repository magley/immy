package handler

import (
	"immy-api/service"
	util "immy-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MetaHandler struct {
	BoardService *service.BoardService
}

func (h *MetaHandler) GetMimeTypes(c *gin.Context) {
	// https://mimetype.io/all-types
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/MIME_types/Common_types
	// https://www.iana.org/assignments/media-types/media-types.xhtml

	mimeTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"video/webm",
		"video/mp4",
	}
	util.OK(c, mimeTypes)
}

func (h *MetaHandler) GetStatistics(c *gin.Context) {
	boardStats, err := h.BoardService.GetStatistics()

	if err != nil {
		util.Fail(c, http.StatusBadRequest, "LIST_FAIL", err.Error())
		return
	} else {
		util.OK(c, boardStats)
		return
	}
}