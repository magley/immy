package route

import (
	"immy-api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterConfigRoutes(h *handler.ConfigHandler, rg *gin.RouterGroup) {
	banappeals := rg.Group("/config")
	{
		banappeals.GET("/", h.GetConfig)
		banappeals.PUT("/post", h.SetPostingEnabled) // ?enabled=0|not0
	}
}
