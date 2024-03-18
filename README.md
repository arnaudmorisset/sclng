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

## Design Decisions and Improvements

### Some cleaning before starting

Before starting to work on the project, I cleaned up a bit.
First, I removed the Docker-based development setup, opting for something lighter.
Running the application is now done using a local installation of Go and Air, allowing fast live reloading.
The inconvenience is that we now need to have requirements installed locally, but I still find it lightweight to install Go and Air rather than relying on Docker.
Then, I changed the whole file structure of the project to be more aligned with standard Go recommendations for server software architecture: https://go.dev/doc/modules/layout#server-project
I also removed Dave Chenney's error management package, as I prefer to follow standards when I can; the project now uses the `errors` package from the standard library.

### Github Client and Configuration

I've isolated HTTP requests toward Github API in a dedicated Github Client package.
The "sensible" information, such as the API key, is fetched by the configuration package from environment variables.
The handler will rely on the GitHub client to fetch the data and only handle parsing incoming parameters and building the resulting view.

### Concurrency

I usually handle concurrency manually by running closures using the `go` keyword, waiting on them using a waiting group, and stacking errors in a channel.
This time, I chose to rely on a library to gain some time.
One thing I usually dislike about Go concurrency API is that it's up to you to write the "glue" code to add the missing lifecycle orchestrator.
The `conc` library by Sourcegraph is such a collection of utility functions that bring a mechanism of ownership over goroutines followed by nice utils functions.
I use it to parse repository information (fetching programming language stats) through their iterator mechanism, and it's pleasant to use.

### Known Limitations & Failure

My approach to this exercise was the wrong one.
I decided to make it the best Go repo I could, but I've tried to bite off more than I could chew here.
My initial goal was to restructure the project, implement the required features, add observability (logs, traces, and metrics), and a CI/CD pipeline to deploy the project to Scalingo.
I was very optimistic regarding the time I could dedicate to this exercise, so everything went well, and the project has the current issues:

- The endpoint doesn't return the information about the last hundred repos but the first hundred ones. I didn't find anything on GitHub REST API (nor with the public repos endpoint or the search endpoint) to obtain information about the most recent repositories.
- The project has no tests. In a real-world scenario, I would have added Unit Tests on the GitHub Client and integration tests on the handlers.
- The only filter is by language. I needed more time to add more.

I'm pretty disappointed by my performance here, but still, it was a fun project and a nice challenge.
Please don't hesitate to share feedback about my code structure and reasoning here.
Thanks!
