package users

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(h *UserHandler, rg *gin.RouterGroup) {
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