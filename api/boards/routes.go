package boards

import (
	"github.com/gin-gonic/gin"
)

func RegisterBoardRoutes(h *BoardHandler, rg *gin.RouterGroup) {
	boards := rg.Group("/boards")
	{
    	boards.GET("/", h.ListBoards)
    	boards.POST("/", h.CreateBoard)
    	boards.GET("/:code", h.GetBoard)
    	boards.PUT("/:code", h.UpdateBoard)
    	boards.DELETE("/:code", h.DeleteBoard)
	}
}