package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sgitwhyd/cangkruan-api/internal/configs"
	"github.com/sgitwhyd/cangkruan-api/internal/handler/memberships"
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
		log.Fatal("Gagal inisiasi config")
	}

	cfg = configs.Get()
	log.Println("config", cfg)

	
	membershipHandler := memberships.NewHandler(r)
	membershipHandler.RegisterRoute()

	r.Run(cfg.Service.Port) // listen and serve on 0.0.0.0:8080
}
