#!/bin/bash

## Install linux
GOOS=linux go build -o chessapi cmd/api/main.go
go build -o chess cmd/chess/main.go

shasum -a 256 chess-v2.0.0-darwin-amd64.tar.gz
zip -r chess-v2.0.0-darwin-amd64.zip chess
