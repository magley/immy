package route

import (
	"github.com/gin-gonic/gin"
	"immy-api/handler"
)

func RegisterUserRoutes(h *handler.UserHandler, rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
    	users.GET("/", h.ListUsers)
    	users.POST("/", h.CreateUser)
    	users.GET("/:id", h.GetUser)
    	users.PUT("/:id", h.UpdateUser)
    	users.DELETE("/:id", h.DeleteUser)
    	users.POST("/login", h.LoginUser)
	}
}