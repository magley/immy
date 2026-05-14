package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterPostRoutes(h *handler.PostHandler, rg *gin.RouterGroup) {
	posts := rg.Group("/posts")
	{
    	posts.GET("/", h.ListPosts)
    	posts.GET("/num/:boardCode/:postNum", h.GetPostByNum)
    	posts.GET("/thread/:threadId", h.GetPostsByThread)
    	posts.POST("/", h.CreatePost)
    	posts.PUT("/:id", h.UpdatePost)
    	posts.DELETE("/:id", h.DeletePost)
	}
}