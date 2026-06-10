package main

import (
	"log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	
	"immy-api/repo"
	"immy-api/service"
	"immy-api/handler"
	"immy-api/route"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "USER", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))	
	
	boardRepo := &repo.BoardRepo{DB: gormDB}
	postRepo := &repo.PostRepo{DB : gormDB}
	threadRepo := &repo.ThreadRepo{DB : gormDB}
	userRepo := &repo.UserRepo{DB : gormDB}
	
	var boardService *service.BoardService
	var postService *service.PostService
	var threadService *service.ThreadService
	var userService *service.UserService
	
	boardService = &service.BoardService{BoardRepo: boardRepo}
	postService = &service.PostService{PostRepo: postRepo, BoardService: boardService, ThreadRepo: threadRepo}
	threadService = &service.ThreadService{ThreadRepo: threadRepo, BoardService: boardService, PostService: postService}
	userService = &service.UserService{UserRepo: userRepo}
	
	postService.ThreadService = threadService
	
	boardHandler := &handler.BoardHandler{BoardService: boardService}
	postHandler := &handler.PostHandler{PostService: postService, UserService: userService}
	threadHandler := &handler.ThreadHandler{ThreadService: threadService, BoardService: boardService, PostService: postService, UserService: userService}
	userHandler := &handler.UserHandler{UserService: userService}

	metaHandler := &handler.MetaHandler{BoardService: boardService}
	
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			route.RegisterMetaRoutes(metaHandler, v1)
			route.RegisterBoardRoutes(boardHandler, v1)
			route.RegisterUserRoutes(userHandler, v1)
			route.RegisterPostRoutes(postHandler, v1)
			route.RegisterThreadRoutes(threadHandler, v1)
		}
	}
	
	router.Run(":8080")
}