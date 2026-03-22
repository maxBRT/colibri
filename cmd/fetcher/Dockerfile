FROM golang:1.26-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /fetcher ./cmd/fetcher



RUN mkdir -p /app/sources
FROM alpine:latest
COPY --from=build /fetcher /fetcher
COPY ./sources/sources.csv /app/sources/sources.csv
CMD ["/fetcher"]

