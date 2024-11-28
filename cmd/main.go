package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/configs"
	"github.com/sgitwhyd/cangkruan-api/internal/handlers"
	"github.com/sgitwhyd/cangkruan-api/internal/repository"
	"github.com/sgitwhyd/cangkruan-api/internal/service"
	"github.com/sgitwhyd/cangkruan-api/pkg/internalsql"
)

// BasePath /api/v1
func main() {
	err := configs.Init(
		configs.WithConfigFolder([]string{"./internal/configs"}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)

	if err != nil {
		log.Fatal("Gagal inisiasi config", err)
		return
	}

	
	cfg := configs.Get()
	log.Println(cfg)

	db, err := internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatal("gagal inisiasi database", err)
		return
	}

	defer db.Close()
	
	log.Println("db connected")

	r := gin.Default()
	route:= r.Group("/api/v1/")
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	userRepo := repository.NewUserRepository(db)
	userActRepo := repository.NewUserActivityRepository(db)

	commentSvc := service.NewCommentService(commentRepo)
	postSvc := service.NewPostService(postRepo, commentRepo, userActRepo)
	authSvc := service.NewAuthService(cfg, userRepo)
	userActSvc := service.NewUserActivityService(userActRepo)

	commenHandler  := handlers.NewCommentHandler(route, commentSvc, postSvc)
	postHandler := handlers.NewPostHandler(route, postSvc)
	membershipHandler := handlers.NewAuthHandler(route, authSvc)
	userActHandler := handlers.NewUserActHandler(route, userActSvc, postSvc)

	commenHandler.RegisterRoute()
	membershipHandler.RegisterRoute()
	postHandler.RegisterRoute()
	userActHandler.RegisterRoute()

	r.Run(cfg.Service.Port) // listen and serve on 0.0.0.0:8080
}
