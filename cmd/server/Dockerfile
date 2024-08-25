# syntax=docker/dockerfile:1
FROM golang:1.22.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /goph-keeper ./cmd/server/main.go

FROM alpine AS server
WORKDIR /
COPY --from=builder /goph-keeper /goph-keeper
WORKDIR /config
COPY ../../config/dev.yaml ./
EXPOSE 8080
ENTRYPOINT ["/goph-keeper"]