package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kontentski/develops-today-task/config"
	"github.com/Kontentski/develops-today-task/pkg/httpserver"
	"github.com/Kontentski/develops-today-task/pkg/logging"
	"github.com/Kontentski/develops-today-task/pkg/postgresql"
	"github.com/gin-gonic/gin"

	"github.com/Kontentski/develops-today-task/internal/api/cat"
	httpController "github.com/Kontentski/develops-today-task/internal/controller/http"
	"github.com/Kontentski/develops-today-task/internal/entity"
	"github.com/Kontentski/develops-today-task/internal/service"
	"github.com/Kontentski/develops-today-task/internal/storage"
)

// Run - initializes and runs application.
func Run(cfg *config.Config) {
	logger := logging.NewZapLogger(cfg.Log.Level)

	postgresql, err := postgresql.NewPostgreSQLGorm(postgresql.Config{
		User:     cfg.PostgreSQL.User,
		Password: cfg.PostgreSQL.Password,
		Host:     cfg.PostgreSQL.Host,
		Database: cfg.PostgreSQL.Database,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to init repository: %w", err))
	}

	// create UUID extension.
	err = postgresql.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create uuid-ossp extension: %w", err))
	}

	err = postgresql.DB.AutoMigrate(
		&entity.SpyCat{},
		&entity.Mission{},
		&entity.Target{},
	)
	if err != nil {
		log.Fatal(fmt.Errorf("automigration failed: %w", err))
	}

	storages := service.Storages{
		SpyCat:  storage.NewSpyCatStorage(postgresql),
		Mission: storage.NewMissionStorage(postgresql),
		Target:  storage.NewTargetStorage(postgresql),
	}

	apis := service.APIs{
		CatAPI: cat.New(&cat.Options{
			Logger: logger,
			Config: cfg,
		}),
	}

	serviceOptions := service.Options{
		Storages: storages,
		APIs:     apis,
		Config:   cfg,
		Logger:   logger,
	}

	services := service.Services{
		SpyCat:  service.NewSpyCatService(serviceOptions, storages.SpyCat),
		Mission: service.NewMissionService(serviceOptions, storages.Mission),
		Target:  service.NewTargetService(serviceOptions, storages.Target),
	}

	httpHandler := gin.New()

	httpController.New(httpController.Options{
		Handler:  httpHandler,
		Services: services,
		Logger:   logger,
		Config:   cfg,
	})

	httpServer := httpserver.New(
		httpHandler,
		httpserver.Port(cfg.HTTP.Port),
		httpserver.ReadTimeout(time.Second*60),
		httpserver.WriteTimeout(time.Second*60),
		httpserver.ShutdownTimeout(time.Second*30),
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())

	case err = <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify", "err", err)
	}

	err = httpServer.Shutdown()
	if err != nil {
		logger.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}
