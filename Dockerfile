FROM golang:latest as builder

WORKDIR /app

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOCACHE=/root/.cache/go-build \
    GOMODCACHE=/go/pkg/mod

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/main

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY fonts /app/fonts

# 必要なディレクトリを作成
RUN mkdir -p /uploads/images \
    && mkdir -p /rss

# 必要に応じて権限を設定
# USER root

CMD ["/app/main"]