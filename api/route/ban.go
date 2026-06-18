package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterBanRoutes(h *handler.BanHandler, rg *gin.RouterGroup) {
	bans := rg.Group("/bans")
	{
    	bans.GET("/", h.ListBans)
    	bans.GET("/admin", h.ListBansForAdmin)
    	bans.GET("/my", h.GetMyBans)
    	bans.POST("/", h.CreateBan)
    	bans.GET("/:id", h.GetBan)
    	bans.GET("/admin/:id", h.GetBanForAdmin)
    	bans.PUT("/:id", h.UpdateBan)
    	bans.DELETE("/:id", h.DeleteBan)
	}
}