package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/configs"
	"github.com/sgitwhyd/cangkruan-api/internal/handlers"
	"github.com/sgitwhyd/cangkruan-api/internal/repository"
	"github.com/sgitwhyd/cangkruan-api/internal/service"
	"github.com/sgitwhyd/cangkruan-api/pkg/formater"
	"github.com/sgitwhyd/cangkruan-api/pkg/internalsql"
)

// BasePath /api/v1
func main() {
	err := configs.Init(
		configs.WithConfigFolder([]string{"./internal/configs"}),
		configs.WithConfigFile("app"),
		configs.WithConfigType("env"),
	)

	if err != nil {
		log.Fatal("Gagal inisiasi config", err)
		return
	}

	
	cfg := configs.Get()
	log.Println(cfg)

	db, err := internalsql.Connect(cfg.Database)
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Adjust to your frontend URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(func(c *gin.Context) {
    c.Next()
    if len(c.Errors) > 0 {
        log.Println("Middleware Error:", c.Errors)
    }
})

	
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

	r.GET("/", func(ctx *gin.Context) {
		response := formater.APIResponse("App is online", 200, "success", nil)
		ctx.JSON(200, response)
	})

	r.Run(cfg.Port) // listen and serve on 0.0.0.0:8080
}
