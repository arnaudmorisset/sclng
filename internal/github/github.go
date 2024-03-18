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
	GetLastHundredRepos(filters Filters) ([]*Repo, error)
}

type GithubClientImpl struct {
	cfg config.GithubConfig
	log logrus.FieldLogger
	clt *req.Client
}

func NewGithubClient(cfg config.GithubConfig, log logrus.FieldLogger) GithubClient {
	return GithubClientImpl{cfg: cfg, log: log, clt: req.C()}
}

func (g GithubClientImpl) GetLastHundredRepos(filters Filters) ([]*Repo, error) {
	var repos []Repo

	g.log.Info("Getting the last 100 public repositories from Github")
	_, err := g.clt.R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		SetBearerAuthToken(g.cfg.Token).
		SetSuccessResult(&repos).
		Get(g.cfg.BaseURL + "/repositories")

	if err != nil {
		return nil, fmt.Errorf("fail to get the repositories: %s", err.Error())
	}

	// Done concurrently
	reposWithLanguages := iter.Map(repos, func(repo *Repo) *Repo {
		languages, err := g.GetRepoLanguages(repo.LanguagesURL)
		if err != nil {
			g.log.Error("Fail to get the languages for the repository ", repo.FullName, ": ", err.Error())
			return repo
		}
		if filters.Language != "" {
			if _, ok := languages[filters.Language]; !ok {
				return nil
			}
		}
		repo.Languages = languages
		return repo
	})

	return cleanReposSlice(reposWithLanguages), nil
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

func cleanReposSlice(repos []*Repo) []*Repo {
	var cleanedRepos []*Repo
	for _, repo := range repos {
		if repo != nil {
			cleanedRepos = append(cleanedRepos, repo)
		}
	}
	return cleanedRepos
}
