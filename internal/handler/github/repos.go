package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	"github.com/arnaudmorisset/sclng/internal/config"
	"github.com/arnaudmorisset/sclng/internal/github"
)

func NewReposHandler(cfg config.Config) handlers.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		log := logger.Get(r.Context())

		resp, err := github.GetLastHundredRepos(cfg.Github)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			e := fmt.Errorf("fail to get the repositories: %s", err.Error())
			log.WithError(e).Error(e.Error())
			return e
		}

		res, err := json.Marshal(resp)
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
}
