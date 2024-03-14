# Canvas for Backend Technical Test at Scalingo

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
- Secrets in env vars

Feedbacks:

- Not please with how errors are handled at the router level
