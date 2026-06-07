package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterThreadRoutes(h *handler.ThreadHandler, rg *gin.RouterGroup) {
	threads := rg.Group("/threads")
	{
    	threads.GET("/", h.ListThreads)
    	threads.GET("/board/:boardCode", h.ListThreadsOfBoard)
    	threads.GET("/board/:boardCode/:num", h.GetThreadByNum)
    	threads.GET("/board/:boardCode/:num/full", h.GetFullThreadByNum)
    	threads.GET("/board/:boardCode/catalog", h.GetThreadsForCatalog)
    	threads.GET("/board/:boardCode/archive", h.GetThreadsForArchive)
    	threads.GET("/board/:boardCode/home", h.GetThreadsForHome)
    	threads.GET("/:id/full", h.GetFullThread)
    	threads.GET("/:id", h.GetThread)
    	threads.POST("/", h.CreateThread)
    	threads.PUT("/:id", h.UpdateThread)
    	threads.DELETE("/:id", h.DeleteThread)
    	threads.PUT("/:id/archive", h.ArchiveThread)
	}
}