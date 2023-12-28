FROM golang:alpine AS build

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

WORKDIR /app

COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN \
    go mod download && \
    go mod tidy && \
    go build -a -o my-epic-app -ldflags='-s -w -extldflags "-static"' ./cmd/main.go

# Main stage

FROM scratch

ARG version="docker"
ARG commit_sha="abcdefg"

# Env for our app to use

ENV PORT=7883
ENV DB_URL="/var/db/prod.db"
ENV ENCRYPTION_KEY="ABCDEFGHIJKLMNOPQRSTUVWXYZ"
ENV VERSION="v0.0.1+${version}"
ENV COMMIT_SHA="${commit_sha}"

# Copy all the Code and stuff to compile everything

WORKDIR /app
COPY assets ./assets
COPY views ./views
COPY --from=build /app/my-epic-app ./

# Exposes port because our program listens on that port
EXPOSE ${PORT}

VOLUME [ "/var/db" ]

ENTRYPOINT [ "./my-epic-app" ]
