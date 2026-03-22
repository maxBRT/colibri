FROM golang:1.26-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /migrate ./cmd/migrate

RUN mkdir -p /app/sources
FROM alpine:latest
COPY --from=build /migrate /migrate
CMD ["/migrate", "up"]

