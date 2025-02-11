package main

import (
	"github.com/ilyakaznacheev/cleanenv"

	"github.com/Kontentski/develops-today-task/config"
	"github.com/Kontentski/develops-today-task/pkg/logging"

	"github.com/Kontentski/develops-today-task/internal/app"
)

func main() {
	logger := logging.NewZapLogger("main")

	var cfg config.Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		logger.Fatal("failed to read env", "err", err)
	}
	logger.Info("read config", "config", cfg)

	app.Run(&cfg)
}
