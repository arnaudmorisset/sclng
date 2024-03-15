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

type Repo struct {
	FullName   string `json:"full_name"`
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
}

type ReposResponse struct {
	Repositories []Repo `json:"repositories"`
}

func NewReposHandler(cfg config.Config) handlers.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		log := logger.Get(r.Context())

		repos, err := github.GetLastHundredRepos(cfg.Github)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			e := fmt.Errorf("fail to get the repositories: %s", err.Error())
			log.WithError(e).Error(e.Error())
			return e
		}

		res, err := toJSON(repos)
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

func toJSON(repos []github.Repo) ([]byte, error) {
	resp := ReposResponse{}
	for _, repo := range repos {
		resp.Repositories = append(resp.Repositories, Repo{
			FullName:   repo.FullName,
			Owner:      repo.Owner.Login,
			Repository: repo.Name,
		})
	}

	return json.Marshal(resp)
}
