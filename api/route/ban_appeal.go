package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterBanAppealRoutes(h *handler.BanAppealHandler, rg *gin.RouterGroup) {
	banappeals := rg.Group("/banappeals")
	{
    	banappeals.GET("/", h.ListBanAppeals)
    	banappeals.POST("/", h.CreateBanAppeal)
    	banappeals.GET("/ban/:id", h.GetBanAppealsOfBan)
    	banappeals.GET("/ban/:id/can", h.CanAppealBan)
    	banappeals.GET("/:id", h.GetBanAppeal)
    	banappeals.PUT("/:id", h.UpdateBanAppeal)
    	banappeals.DELETE("/:id", h.DeleteBanAppeal)
	}
}