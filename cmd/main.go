package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/configs"
	"github.com/sgitwhyd/cangkruan-api/internal/handler/memberships"
	membershipRepo "github.com/sgitwhyd/cangkruan-api/internal/repository/memberships"
	"github.com/sgitwhyd/cangkruan-api/pkg/internalsql"
)

func main() {
	r := gin.Default()

	var cfg *configs.Config

	err := configs.Init(
		configs.WithConfigFolder([]string{"./internal/configs"}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)

	if err != nil {
		log.Fatal("Gagal inisiasi config", err)
	}

	cfg = configs.Get()
	log.Println(cfg)


	db, err := internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatal("gagal inisiasi database", err)
	}

	_ = membershipRepo.NewRepository(db)

	
	membershipHandler := memberships.NewHandler(r)
	membershipHandler.RegisterRoute()

	r.Run(cfg.Service.Port) // listen and serve on 0.0.0.0:8080
}
