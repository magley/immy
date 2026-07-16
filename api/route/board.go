package route

import (
	"immy-api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterBoardRoutes(h *handler.BoardHandler, rg *gin.RouterGroup) {
	boards := rg.Group("/boards")
	{
		boards.GET("/", h.GetAllBoards)
		boards.POST("/", h.CreateBoard)
		boards.GET("/id/:id", h.GetBoardById)
		boards.GET("/code/:code", h.GetBoard)
		boards.PUT("/:code", h.UpdateBoard)
		boards.DELETE("/:code", h.DeleteBoard)
	}
}
