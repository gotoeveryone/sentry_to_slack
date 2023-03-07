# Sentry to Slack

Sentry のデータを Slack に連携するためのツール

## Requirements

- Docker

## Run

```console
$ cp .env.example .env
$ docker compose up
```

## Post data

```console
$ curl -X POST -d '{"url":"https://test.example.com","event":{"title":"TestError","level":"error","environment": "test"}}' http://localhost:8080
```
