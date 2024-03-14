package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arnaudmorisset/sclng/internal/config"
	"github.com/arnaudmorisset/sclng/internal/github"
)

func TestGetLastHundredReposReturnsRepos(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"full_name": "repo1", "owner": {"login": "owner1"}}, {"full_name": "repo2", "owner": {"login": "owner2"}}]`)
	}))
	defer mockServer.Close()

	cfg := config.Config{
		Github: config.GithubConfig{
			Token:   "<YOUR_TOKEN>",
			BaseURL: mockServer.URL,
		},
	}

	expectedRepos := []github.Repo{
		{FullName: "repo1", Owner: github.Owner{Login: "owner1"}},
		{FullName: "repo2", Owner: github.Owner{Login: "owner2"}},
	}

	handler := NewReposHandler(cfg)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	err := handler(rec, req, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var actualRepos []github.Repo
	err = json.NewDecoder(res.Body).Decode(&actualRepos)
	if err != nil {
		t.Fatalf("failed to decode response body: %s", err)
	}

	if len(actualRepos) != len(expectedRepos) {
		t.Errorf("expected %d repos, got %d", len(expectedRepos), len(actualRepos))
	}

	for i, expectedRepo := range expectedRepos {
		actualRepo := actualRepos[i]
		if actualRepo.FullName != expectedRepo.FullName {
			t.Errorf("expected repo name to be '%s', got '%s'", expectedRepo.FullName, actualRepo.FullName)
		}
		if actualRepo.Owner.Login != expectedRepo.Owner.Login {
			t.Errorf("expected repo owner to be '%s', got '%s'", expectedRepo.Owner.Login, actualRepo.Owner.Login)
		}
	}
}
