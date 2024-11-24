package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/configs"
	"github.com/sgitwhyd/cangkruan-api/internal/handlers"
	"github.com/sgitwhyd/cangkruan-api/internal/handlers/memberships"
	"github.com/sgitwhyd/cangkruan-api/internal/repository"
	membershipRepo "github.com/sgitwhyd/cangkruan-api/internal/repository/memberships"
	"github.com/sgitwhyd/cangkruan-api/internal/service"
	membershipSvc "github.com/sgitwhyd/cangkruan-api/internal/service/memberships"
	"github.com/sgitwhyd/cangkruan-api/pkg/internalsql"
)

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
	r.Use(gin.Logger())
	r.Use(gin.Recovery())



	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	membershipRepository := membershipRepo.NewRepository(db)

	commentScv := service.NewCommentService(commentRepo)
	postScv := service.NewPostService(postRepo)
	membershipSvc := membershipSvc.NewService(cfg, membershipRepository)

	commenHandler  := handlers.NewCommentHandler(r, commentScv, postScv)
	postHandler := handlers.NewPostHandler(r, postScv)
	membershipHandler := memberships.NewHandler(r, membershipSvc)

	commenHandler.RegisterRoute()
	membershipHandler.RegisterRoute()
	postHandler.RegisterRoute()

	r.Run(cfg.Service.Port) // listen and serve on 0.0.0.0:8080
}
