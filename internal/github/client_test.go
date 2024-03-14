package github

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arnaudmorisset/sclng/internal/config"
)

func TestGetLastHundredRepos(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"full_name": "repo1", "owner": {"login": "owner1"}}, {"full_name": "repo2", "owner": {"login": "owner2"}}]`)
	}))
	defer mockServer.Close()

	cfg := config.GithubConfig{
		Token:   "<YOUR_TOKEN>",
		BaseURL: mockServer.URL,
	}

	repos, err := GetLastHundredRepos(cfg)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}

	if len(repos) != 2 {
		t.Errorf("expected 2 repos, got %d", len(repos))
	}
	if repos[0].FullName != "repo1" {
		t.Errorf("expected first repo name to be 'repo1', got '%s'", repos[0].FullName)
	}
	if repos[0].Owner.Login != "owner1" {
		t.Errorf("expected first repo owner to be 'owner1', got '%s'", repos[0].Owner.Login)
	}
	if repos[1].FullName != "repo2" {
		t.Errorf("expected second repo name to be 'repo2', got '%s'", repos[1].FullName)
	}
	if repos[1].Owner.Login != "owner2" {
		t.Errorf("expected second repo owner to be 'owner2', got '%s'", repos[1].Owner.Login)
	}
}
