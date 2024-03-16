package github

import (
	"fmt"

	"github.com/arnaudmorisset/sclng/internal/config"
	"github.com/imroc/req/v3"
	"github.com/sourcegraph/conc/iter"
)

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
}

type GithubClientImpl struct {
	cfg config.GithubConfig
}

func NewGithubClient(cfg config.GithubConfig) GithubClient {
	return GithubClientImpl{cfg: cfg}
}

func (g GithubClientImpl) GetLastHundredRepos() ([]Repo, error) {
	var repos []Repo

	c := req.C()

	_, err := c.R().
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
		languages, err := g.GetRepoLanguages(repo.LanguagesURL)
		if err != nil {
			return
		}
		repo.Languages = languages
	})

	return repos, nil
}

func (g GithubClientImpl) GetRepoLanguages(languagesURL string) (map[string]int, error) {
	var languages map[string]int

	c := req.C()

	_, err := c.R().
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
