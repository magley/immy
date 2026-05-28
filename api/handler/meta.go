package handler

import (
	"github.com/gin-gonic/gin"
	util "immy-api/util"
)

type MetaHandler struct {}

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