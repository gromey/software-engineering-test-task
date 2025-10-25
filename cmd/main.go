package main

import (
	"log"

	"cruder/internal/config"
	"cruder/internal/controller"
	"cruder/internal/handler"
	"cruder/internal/repository"
	"cruder/internal/service"
	"cruder/pkg/logger"

	"github.com/easysy/envio"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := new(config.Config)
	if err := envio.Get(cfg); err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	logger.SetLogger(cfg.LogLevel)

	dbConn, err := repository.NewPostgresConnection(cfg.GetPostgresDNS())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	repositories := repository.NewRepository(dbConn)
	services := service.NewService(repositories)
	controllers := controller.NewController(services)

	r := gin.Default()
	handler.New(r, cfg.APIKey, controllers.Users)

	if err = r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
