# chess
[![CircleCI](https://circleci.com/gh/dumbogo/chess.svg?style=shield)](https://circleci.com/gh/dumbogo/chess)
[![Go Report Card](https://goreportcard.com/badge/github.com/dumbogo/chess)](https://goreportcard.com/report/github.com/dumbogo/chess)
[![GoDoc](https://pkg.go.dev/badge/github.com/dumbogo/chess?status.svg)](https://pkg.go.dev/github.com/dumbogo/chess)

a chess p2p game using CLI

## Server
### Install

#### Building from source
```sh
$ go get github.com/dumbogo/chess/cmd/chessapi
$ chessapi help

```
#### Dependencies
In order to run the server, you need to install some services:
- Postgresql
- NATS

You can use docker-compose:
```sh
$ docker network create chess
$ docker-compose up -d
```

Or you can instead use `deps.sh`:
```sh
$ sh deps.sh
```

### Run
In order to run the server you need to create a `config.toml` file first, sample:
```TOML
ENV= "development"

[API]
port = ":8000"
server_cert = "/opt/data/chessapi/x509/server_cert.pem"
server_key = "/opt/data/chessapi/x509/server_key.pem"

[Database]
host = "127.0.0.1"
port = "5432"
db_name = "chess_api"

[HTTP_server]
Scheme = "http"
# This is necessary to configure in you auth github callback configuration
Host = "yourdomainorip.com"
Port = ":8080"
```
In order to be able to use github auth, you need to configure a github application and oauth2

Make sure you have the corresponding `server_cert` and `server_key` on your system, the repository has some pregenerated files within `certs` directory.
WARNING! its only for dev purposes

Also you need to define some env variables:
```bash
export CHESS_API_DATABASE_USERNAME=postgres
export CHESS_API_DATABASE_PASSWORD=password

export CHESS_API_GITHUB_KEY=key
export CHESS_API_GITHUB_SECRET=secret

export CHESS_API_NATS_URL=localhost:4222
```


Once you have everything in place, you need to run migrations first, then start the server:
```sh
$ chessapi migrate -c config.toml
$ chessapi start -c config.toml
```

## Client
### Install

#### MacOSX
With brew:
```sh
$ brew install dumbogo/tap/chess
```

### Building from source
```sh
$ go get github.com/dumbogo/chess/cmd/chess
$ chess help
```

### Configure
To configure your client, you need to add a TOML config file on `$HOME/.chess/config`, you can use default config from command bellow:
```sh
# create folders:
$ mkdir -p ~/.chess/certs/x509
$ chess config default > ~/.chess/config
```

Also, you need to add the client certfile on `$HOME/.chess/certs/` location, you can use the samples on `certs`.
WARNING! the certs on `certs` folder are for dev purpoeses only

### Play
```sh
$ chess help
Chess multi-player game on terminal

Usage:
  chess [command]

Available Commands:
  help        Help about any command
  join        Join game
  move        Move piece
  signup      Sign up on chess
  start       start game
  version     Print chess version
  watch       watch game

Flags:
  -h, --help   help for chess

Use "chess [command] --help" for more information about a command.
```
