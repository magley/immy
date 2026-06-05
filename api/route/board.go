package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterBoardRoutes(h *handler.BoardHandler, rg *gin.RouterGroup) {
	boards := rg.Group("/boards")
	{
    	boards.GET("/", h.ListBoards)
    	boards.POST("/", h.CreateBoard)
    	boards.GET("/id/:id", h.GetBoardById)
    	boards.GET("/code/:code", h.GetBoard)
    	boards.PUT("/:code", h.UpdateBoard)
    	boards.DELETE("/:code", h.DeleteBoard)
	}
}