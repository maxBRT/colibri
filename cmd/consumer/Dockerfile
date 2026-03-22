FROM golang:1.26-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /consumer ./cmd/consumer

RUN mkdir -p /app/config

FROM alpine:latest
COPY --from=build /consumer /consumer
COPY ./config/prompts/description_agent_prompt.md /app/config/description_agent_prompt.md
CMD ["/consumer"]
