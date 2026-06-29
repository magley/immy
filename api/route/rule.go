package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterRuleRoutes(h *handler.RuleHandler, rg *gin.RouterGroup) {
	Rules := rg.Group("/rules")
	{
    	Rules.POST("/board/", h.CreateRuleBoard)
    	Rules.DELETE("/board/:boardId/:ruleId", h.DeleteRule)
    	// TODO: These three could be a single endpoint with [optional] query params
    	Rules.GET("/board/", h.ListAllRuleBoards)
    	Rules.GET("/board/rules/:boardId", h.ListAllRulesOfBoard)
    	Rules.GET("/board/boards/:ruleId", h.ListAllBoardsOfRule)

    	Rules.GET("/", h.ListRules)
    	Rules.POST("/", h.CreateRule)
    	Rules.GET("/:id", h.GetRule)
    	Rules.PUT("/:id", h.UpdateRule)
    	Rules.DELETE("/:id", h.DeleteRule)
	}
}