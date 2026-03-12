#!/bin/bash

start_or_run () {
    docker inspect rabbitmq_colibri > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "Starting Colibri RabbitMQ container..."
        docker start rabbitmq_colibri
    else
        echo "Colibri RabbitMQ container not found, creating a new one..."
        docker run -d --name rabbitmq_colibri -p 5672:5672 -p 15672:15672 -p 61613 rabbitmq-stomp
    fi
}

case "$1" in
    start)
        start_or_run
        ;;
    stop)
        echo "Stopping Colibri RabbitMQ container..."
        docker stop rabbitmq_colibri
        ;;
    logs)
        echo "Fetching logs for Colibri RabbitMQ container..."
        docker logs -f rabbitmq_colibri
        ;;
    *)
        echo "Usage: $0 {start|stop|logs}"
        exit 1
esac
