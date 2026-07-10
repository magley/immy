package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"immy-api/handler"
	"immy-api/repo"
	"immy-api/route"
	"immy-api/service"
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_DBNAME"), os.Getenv("DB_PASSWORD"),
	)
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
	postRepo := &repo.PostRepo{DB: gormDB}
	threadRepo := &repo.ThreadRepo{DB: gormDB}
	userRepo := &repo.UserRepo{DB: gormDB}
	banRepo := &repo.BanRepo{DB: gormDB}
	banAppealRepo := &repo.BanAppealRepo{DB: gormDB}
	blogpostRepo := &repo.BlogpostRepo{DB: gormDB}
	ruleRepo := &repo.RuleRepo{DB: gormDB}

	var boardService *service.BoardService
	var postService *service.PostService
	var threadService *service.ThreadService
	var userService *service.UserService
	var banService *service.BanService
	var banAppealService *service.BanAppealService
	var blogpostService *service.BlogpostService
	var ruleService *service.RuleService

	boardService = &service.BoardService{BoardRepo: boardRepo}
	postService = &service.PostService{PostRepo: postRepo, BoardService: boardService, ThreadRepo: threadRepo}
	threadService = &service.ThreadService{ThreadRepo: threadRepo, BoardService: boardService, PostService: postService}
	userService = &service.UserService{UserRepo: userRepo}
	banService = &service.BanService{BanRepo: banRepo, BoardService: boardService, UserService: userService}
	banAppealService = &service.BanAppealService{BanAppealRepo: banAppealRepo, BanService: banService, UserService: userService}
	blogpostService = &service.BlogpostService{BlogpostRepo: blogpostRepo, UserService: userService}
	ruleService = &service.RuleService{RuleRepo: ruleRepo}

	postService.ThreadService = threadService

	boardHandler := &handler.BoardHandler{BoardService: boardService}
	postHandler := &handler.PostHandler{PostService: postService, UserService: userService, BanService: banService}
	threadHandler := &handler.ThreadHandler{ThreadService: threadService, BoardService: boardService, PostService: postService, UserService: userService, BanService: banService}
	userHandler := &handler.UserHandler{UserService: userService}
	banHandler := &handler.BanHandler{BanService: banService, UserService: userService}
	banAppealHandler := &handler.BanAppealHandler{BanAppealService: banAppealService, UserService: userService}
	blogpostHandler := &handler.BlogpostHandler{BlogpostService: blogpostService, UserService: userService}
	ruleHandler := &handler.RuleHandler{RuleService: ruleService}

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
			route.RegisterBanRoutes(banHandler, v1)
			route.RegisterBanAppealRoutes(banAppealHandler, v1)
			route.RegisterBlogpostRoutes(blogpostHandler, v1)
			route.RegisterRuleRoutes(ruleHandler, v1)
		}
	}

	router.Run(":8080")
}
