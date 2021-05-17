#!/bin/bash

GOOS=linux go build -o chessapi cmd/api/main.go
go build -o chess cmd/client/main.go
