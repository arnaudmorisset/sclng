package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	"github.com/arnaudmorisset/sclng/internal/config"
	"github.com/arnaudmorisset/sclng/internal/handler/github"
	"github.com/arnaudmorisset/sclng/internal/handler/healthcheck"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logger.Default()
	log.Info("starting up the application")

	if err := run(log); err != nil {
		log.WithError(err).Error("error running the application")
		os.Exit(1)
	}
}

func run(log logrus.FieldLogger) error {
	log.Info("parsing the configuration")
	cfg, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("fail to parse the configuration: %s", err.Error())
	}

	log.Info(("initializing API"))
	router := handlers.NewRouter(log)
	router.HandleFunc("/ping", healthcheck.NewPongHandler())
	router.HandleFunc("/repos", github.NewReposHandler(cfg))

	log.WithField("port", cfg.Port).Info("listening...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		return fmt.Errorf("fail to listen to the given port: %s", err.Error())
	}

	return nil
}
