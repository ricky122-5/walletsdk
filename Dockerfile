# syntax=docker/dockerfile:1.6

ARG GO_VERSION=1.25.2

FROM golang:${GO_VERSION}-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -trimpath -ldflags="-s -w" -o /out/walletsdk ./cmd/server

FROM gcr.io/distroless/base-debian12

ENV APP_PORT=8080

COPY --from=builder /out/walletsdk /usr/local/bin/walletsdk

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/walletsdk"]

