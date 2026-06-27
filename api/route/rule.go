package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterRuleRoutes(h *handler.RuleHandler, rg *gin.RouterGroup) {
	Rules := rg.Group("/Rules")
	{
    	Rules.GET("/", h.ListRules)
    	Rules.POST("/", h.CreateRule)
    	Rules.GET("/:id", h.GetRule)
    	Rules.PUT("/:id", h.UpdateRule)
    	Rules.DELETE("/:id", h.DeleteRule)
	}
}