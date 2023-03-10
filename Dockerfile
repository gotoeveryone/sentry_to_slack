FROM golang:1.19-alpine as development

ENV LANG C.UTF-8
ENV APP_ROOT /var/app

# hadolint ignore=DL3018
RUN apk add gcc g++ --no-cache

RUN go install github.com/cosmtrek/air@v1.29.0 && \
  go install honnef.co/go/tools/cmd/staticcheck@2022.1.2

WORKDIR ${APP_ROOT}
COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

FROM golang:1.19-alpine as builder

ENV LANG C.UTF-8
ENV APP_ROOT /var/app

# hadolint ignore=DL3018
RUN apk add gcc g++ --no-cache

WORKDIR ${APP_ROOT}
COPY ./ ${APP_ROOT}
RUN go mod download && \
  go build -ldflags '-s -w' -o sentry_to_slack ${APP_ROOT}/src/cmd/sentry_to_slack/main.go

FROM golang:1.19-alpine as production

ENV LANG C.UTF-8
ENV APP_ROOT /var/app

WORKDIR ${APP_ROOT}
COPY --from=builder ${APP_ROOT}/sentry_to_slack ${APP_ROOT}

CMD ["./sentry_to_slack"]
