package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	"github.com/arnaudmorisset/sclng/internal/config"
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
	router.HandleFunc("/ping", pongHandler)

	// Initialize web server and configure the following routes:
	// GET /repos
	// GET /stats

	log.WithField("port", cfg.Port).Info("listening...")
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
	if err != nil {
		return fmt.Errorf("fail to listen to the given port: %s", err.Error())
	}

	return nil
}

func pongHandler(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())

	res, err := json.Marshal(map[string]string{"status": "pong"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e := fmt.Errorf("fail to encode JSON: %s", err.Error())
		log.WithError(e).Error(e.Error())
		return e
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(res)

	return nil
}
