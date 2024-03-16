package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Scalingo/go-handlers"
	"github.com/Scalingo/go-utils/logger"
	"github.com/arnaudmorisset/sclng/internal/config"
	"github.com/arnaudmorisset/sclng/internal/github"
)

type LanguageStats struct {
	Bytes int `json:"bytes"`
}

type Repo struct {
	FullName   string                   `json:"full_name"`
	Owner      string                   `json:"owner"`
	Repository string                   `json:"repository"`
	Languages  map[string]LanguageStats `json:"languages"`
}

type ReposResponse struct {
	Repositories []Repo `json:"repositories"`
}

func NewReposHandler(cfg config.Config, gh github.GithubClient) handlers.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		log := logger.Get(r.Context())

		params, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			e := fmt.Errorf("fail to parse the query parameters: %s", err.Error())
			log.WithError(e).Error(e.Error())
			return e
		}

		filters := github.Filters{
			Language: params.Get("language"),
			License:  params.Get("license"),
		}

		var repos []github.Repo
		if filters.Language != "" || filters.License != "" {
			repos, err = gh.GetLastHundredReposFiltered(filters)
		} else {
			repos, err = gh.GetLastHundredRepos()
		}

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
		languageStats := make(map[string]LanguageStats)
		for lang, bytes := range repo.Languages {
			languageStats[lang] = LanguageStats{Bytes: bytes}
		}

		resp.Repositories = append(resp.Repositories, Repo{
			FullName:   repo.FullName,
			Owner:      repo.Owner.Login,
			Repository: repo.Name,
			Languages:  languageStats,
		})
	}

	return json.Marshal(resp)
}
