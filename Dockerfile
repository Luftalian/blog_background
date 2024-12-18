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

COPY uploads/images/d2f9e27e-aff4-4ce2-ba22-7fd1e614c4f6_thumb.png /uploads/images/d2f9e27e-aff4-4ce2-ba22-7fd1e614c4f6_thumb.png

COPY rss/rss.xml /rss/rss.xml

RUN mkdir -p /uploads/images

RUN mkdir /rss

USER root

CMD ["/app/main"]