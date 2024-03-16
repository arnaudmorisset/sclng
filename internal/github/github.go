package github

import (
	"fmt"

	"github.com/arnaudmorisset/sclng/internal/config"
	"github.com/imroc/req/v3"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/iter"
)

type Filters struct {
	Language string
	License  string
}

type Owner struct {
	Login string `json:"login"`
}

type Repo struct {
	FullName     string         `json:"full_name"`
	Owner        Owner          `json:"owner"`
	Name         string         `json:"name"`
	LanguagesURL string         `json:"languages_url"`
	Languages    map[string]int `json:"languages"`
}

type GithubClient interface {
	GetLastHundredRepos() ([]Repo, error)
	GetLastHundredReposFiltered(filters Filters) ([]Repo, error)
}

type GithubClientImpl struct {
	cfg config.GithubConfig
	log logrus.FieldLogger
	clt *req.Client
}

func NewGithubClient(cfg config.GithubConfig, log logrus.FieldLogger) GithubClient {
	return GithubClientImpl{cfg: cfg, log: log, clt: req.C()}
}

// FIX: this method returns the first 100 repositories, not the last 100
func (g GithubClientImpl) GetLastHundredRepos() ([]Repo, error) {
	var repos []Repo

	g.log.Info("Getting the last 100 public repositories from Github")
	_, err := g.clt.R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		SetBearerAuthToken(g.cfg.Token).
		SetSuccessResult(&repos).
		Get(g.cfg.BaseURL + "/repositories")

	if err != nil {
		return repos, fmt.Errorf("fail to get the repositories: %s", err.Error())
	}

	// Done concurrently
	iter.ForEach(repos, func(repo *Repo) {
		g.log.Info("Getting the languages for the repository ", repo.FullName)
		languages, err := g.GetRepoLanguages(repo.LanguagesURL)
		if err != nil {
			g.log.Error("Fail to get the languages for the repository ", repo.FullName, ": ", err.Error())
			return
		}
		repo.Languages = languages
	})

	return repos, nil
}

func (g GithubClientImpl) GetLastHundredReposFiltered(filters Filters) ([]Repo, error) {
	repos, err := g.GetLastHundredRepos()
	if err != nil {
		return repos, err
	}

	// TODO: filter the repositories based on the filters

	return repos, nil
}

func (g GithubClientImpl) GetRepoLanguages(languagesURL string) (map[string]int, error) {
	var languages map[string]int

	g.log.Info("Getting the languages from ", languagesURL)
	_, err := g.clt.R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		SetBearerAuthToken(g.cfg.Token).
		SetSuccessResult(&languages).
		Get(languagesURL)

	if err != nil {
		return languages, fmt.Errorf("fail to get the languages for the repository %s: %s", languagesURL, err.Error())
	}

	return languages, nil
}
