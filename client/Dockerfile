# syntax=docker/dockerfile:1
FROM golang:1.22.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY .. ./
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /goph-keeper-client ./client/main.go

FROM alpine AS client
WORKDIR /
COPY --from=builder /goph-keeper-client /goph-keeper-client
WORKDIR /config
COPY ../../client/config/dev.yaml ./
EXPOSE 8081
ENTRYPOINT ["/goph-keeper-client"]