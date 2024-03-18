# SCLNG - Simple Go HTTP Server for Listing GitHub Repositories

SCLNG is a lightweight Go application designed to run an HTTP server that fetches and lists GitHub repositories.
With minimal configuration, it provides a simple endpoint for retrieving information about public repositories, including their names, owners, and programming languages.

## Features

- HTTP Server: SCLNG starts an HTTP server on port 8080, ready to receive incoming requests.
- `/repos` Endpoint: The server provides a single endpoint /repos to list GitHub repositories.
- GitHub API Integration: Utilizes GitHub REST API to fetch repository information.
- Filtering by Language: Users can specify a programming language through query parameters to filter repositories.

## Usage

```bash
# Create a local-only version of the env vars file
cp .env_vars.sample .env_vars

# Replace the value of the variable GITHUB_TOKEN with your own token
```

```bash
# Run the application (using Go run)
go run cmd/maing.go

# For live-reload: Install air, then run it from the root folder
go install github.com/cosmtrek/air@latest
air
```

```bash
# Fetching repos data (using cURL)
curl "http://localhost:8080/repos"

# Filtering by language (e.g. with Go)
curl "http://localhost:8080/repos?language=go"
```

### Example response

```json
[
  {
    "full_name": "collectiveidea/audited",
    "owner": "collectiveidea",
    "repository": "audited",
    "languages": {
      "JavaScript": {
        "bytes": 49
      },
      "Ruby": {
        "bytes": 134218
      }
    }
  },
  {
    "full_name": "collectiveidea/calendar_builder",
    "owner": "collectiveidea",
    "repository": "calendar_builder",
    "languages": {
      "JavaScript": {
        "bytes": 20126
      },
      "Ruby": {
        "bytes": 34700
      }
    }
  }
]
```

<!-- # Canvas for Backend Technical Test at Scalingo

## Instructions

- From this canvas, respond to the project which has been communicated to you by our team
- Feel free to change everything

## Execution

```
docker compose up
```

Application will be then running on port `5000`

## Test

```
$ curl localhost:5000/ping
{ "status": "pong" }
```

## Notes

> Things I should talk about to explain my tech choices.

- Moving to architecture recommended for Go servers
- Removing docker for local env and using only air for live reload
- Using stdlib for errors instead of Dave Chenney's package
- Reworking main function
- Secrets in env vars -->
