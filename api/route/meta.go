package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterMetaRoutes(h *handler.MetaHandler, rg *gin.RouterGroup) {
	posts := rg.Group("/meta")
	{
    	posts.GET("/mime", h.GetMimeTypes)
	}
}