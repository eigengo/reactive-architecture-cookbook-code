#!/bin/sh
IP=`ifconfig | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p'`
# docker run -d -p 2181:2181 -p 9092:9092 --env ADVERTISED_HOST=$IP --env ADVERTISED_PORT=9092 spotify/kafka
# docker run -d --name summary-redis -p 6379:6379 redis
