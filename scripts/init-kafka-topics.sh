#!/bin/sh
set -e

echo "Waiting for kafka"

until /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --list > /dev/null 2>&1; do
    sleep 2
done

echo "Creating topics"

# Топик для событий "алерт сработал"
/opt/kafka/bin/kafka-topics.sh --create \
  --bootstrap-server kafka:9092 \
  --topic alert.triggered \
  --partitions 3 \
  --replication-factor 1 \
  --if-not-exists

# Топик для событий "уведомление отправлено" (пока не юзаем, логирование, аудит)
/opt/kafka/bin/kafka-topics.sh --create \
  --bootstrap-server kafka:9092 \
  --topic notification.sent \
  --partitions 3 \
  --replication-factor 1 \
  --if-not-exists

echo "Topics:"
/opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --list