FROM golang:1.24-alpine AS build
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/api ./cmd/api

FROM alpine:3.20 AS api
WORKDIR /app
COPY --from=build /app/bin/api /app/api
EXPOSE 8080
CMD ["/app/api"]

FROM golang:1.24-alpine AS goose
WORKDIR /app
RUN apk add --no-cache ca-certificates
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.21.1
ENTRYPOINT ["/go/bin/goose"]
