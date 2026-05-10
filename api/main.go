package main

import (
	"log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	boards "immy-api/boards"
	users "immy-api/users"
)


func main() {
	dsn := "host=db port=5432 user=admin dbname=admin password=admin sslmode=disable"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		
		log.Fatal(err)
	}
	db, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	router := gin.Default()
	router.Use(cors.Default())
	
	boardRepo := &boards.BoardRepo{DB: gormDB}
	boardHandler := &boards.BoardHandler{BoardRepo: boardRepo}
	
	userRepo := &users.UserRepo{DB : gormDB}
	userHandler := &users.UserHandler{UserRepo: userRepo}
	
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			boards.RegisterBoardRoutes(boardHandler, v1)
			users.RegisterUserRoutes(userHandler, v1)
		}
	}
	
	router.Run(":8080")
}
			