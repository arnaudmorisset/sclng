package github

import (
	"fmt"

	"github.com/arnaudmorisset/sclng/internal/config"
	"github.com/imroc/req/v3"
)

type Owner struct {
	Login string `json:"login"`
}

type Repo struct {
	FullName string `json:"full_name"`
	Owner    Owner  `json:"owner"`
	Name     string `json:"name"`
}

func GetLastHundredRepos(cfg config.GithubConfig) ([]Repo, error) {
	var resp []Repo

	c := req.C()

	_, err := c.R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		SetBearerAuthToken(cfg.Token).
		SetSuccessResult(&resp).
		Get(cfg.BaseURL + "/repositories")

	if err != nil {
		return resp, fmt.Errorf("fail to get the repositories: %s", err.Error())
	}

	return resp, nil
}
