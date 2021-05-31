#!/bin/bash

docker network create chess || true
docker run -d \
	--name chess_postgresql \
	--network chess \
	-e POSTGRES_PASSWORD=password \
	-e POSTGRES_DB=chess_api \
	-p 5432:5432 \
	postgres

docker run -d \
	--name chess_nats \
	--network chess \
	-p 4222:4222 \
	nats
