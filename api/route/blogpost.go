package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterBlogpostRoutes(h *handler.BlogpostHandler, rg *gin.RouterGroup) {
	blogposts := rg.Group("/blogposts")
	{
    	blogposts.GET("/", h.ListBlogposts)
    	blogposts.POST("/", h.CreateBlogpost)
    	blogposts.GET("/:id", h.GetBlogpost)
    	blogposts.PUT("/:id", h.UpdateBlogpost)
    	blogposts.DELETE("/:id", h.DeleteBlogpost)
	}
}