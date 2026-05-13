package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterThreadRoutes(h *handler.ThreadHandler, rg *gin.RouterGroup) {
	threads := rg.Group("/threads")
	{
    	threads.GET("/", h.ListThreads)
    	threads.GET("/board/:boardCode", h.ListThreads)
    	threads.POST("/", h.CreateThread)
    	threads.GET("/:id", h.GetThread)
    	threads.PUT("/:id", h.UpdateThread)
    	threads.DELETE("/:id", h.DeleteThread)
	}
}