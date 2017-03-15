#!/bin/sh

# This script builds an integration test Cassandra image seeded with data/init.cql

CASSANDRA_TAG=eas-it-cassandra
CASSANDRA_NAME=eas-it-cassandra
docker build -f Dockerfile.it . --tag $CASSANDRA_TAG

# For convenience, remove old containers
docker ps -f "ancestor=eas-it-cassandra" -qa | xargs docker rm -fv

# Start the built cassandra container
docker run --name $CASSANDRA_NAME -p 9042:9042 -d $CASSANDRA_TAG

# UGH!!
sleep 30

# Initialize it with the data
docker run -it --link $CASSANDRA_NAME:cassandra --rm $CASSANDRA_TAG sh -c 'exec cqlsh "$CASSANDRA_PORT_9042_TCP_ADDR" -f /data/init.cql'

# Run cqlsh
docker run -it --link $CASSANDRA_NAME:cassandra --rm $CASSANDRA_TAG sh -c 'exec cqlsh "$CASSANDRA_PORT_9042_TCP_ADDR"'
