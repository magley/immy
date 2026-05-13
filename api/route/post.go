package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterPostRoutes(h *handler.PostHandler, rg *gin.RouterGroup) {
	// posts := rg.Group("/posts")
	// {
    // 	posts.GET("/", h.ListPosts)
    // 	posts.POST("/", h.CreatePost)
    // 	posts.GET("/:code", h.GetPost)
    // 	posts.PUT("/:code", h.UpdatePost)
    // 	posts.DELETE("/:code", h.DeletePost)
	// }
}