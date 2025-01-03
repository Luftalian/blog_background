FROM golang:latest as builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/go/pkg/mod

COPY go.mod go.sum ./
RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=${GOMODCACHE} \
    go mod download

COPY . .

RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=${GOMODCACHE} \
    go build -o /app/main

FROM gcr.io/distroless/static-debian11:latest

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY fonts /app/fonts
# COPY service-account.json /app/service-account.json

# COPY rss/rss.xml /rss/rss.xml

USER root

CMD ["/app/main"]