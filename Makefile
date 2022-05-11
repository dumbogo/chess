# Go command to use for build
GO ?= go1.16

BINDIR ?= bin
RELEASEDIR ?= releases

## HARDCODED, needs to automatically detects current version by git
version ?= v1.0.0-alpha.3

.PHONY: test integration release build clean proto deployk8 deps deps-clean

test: # run unit tests
	 $(GO) test ./... -cover -coverprofile=coverage.out -v

integration: # run unit + integration tests
	 $(GO) test ./... -cover -coverprofile=coverage.out -tags=integration -v

release: build # release chessapi and chess, only support for Linux chessapi and MacOSX chess client
	cd $(BINDIR) && chmod 755 chessapi chess && \
		tar -czf chessapi-$(version).linux-amd64.tar.gz chessapi && \
		shasum -a 256 chessapi-$(version).linux-amd64.tar.gz > ../releases/sha256sums.txt && \
		mv chessapi-$(version).linux-amd64.tar.gz ../releases/ && \
		zip -r chess-$(version).darwin-amd64.zip chess && \
		shasum -a 256 chess-$(version).darwin-amd64.zip >> ../releases/sha256sums.txt && \
		mv chess-$(version).darwin-amd64.zip ../releases/

build: # Build chess for MacOSX & chessapi for Linux
	$(GO) build -o $(BINDIR)/chess cmd/chess/main.go
	GOOS=linux $(GO) build -o $(BINDIR)/chessapi cmd/chessapi/main.go

clean:
	rm -rf bin/* && rm -rf releases/*.tar.gz releases/sha256sums.txt releases/*.zip

proto: # Build proto files
	GO111MODULE=on protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/service.proto
deployk8:
	echo "starting kubernetes deployment" && \
		kubectl apply -f k8/secrets.yaml && \
		kubectl apply -f k8/configmaps.yml && \
		kubectl apply -f k8/postgresql/deployment.yml && \
		kubectl apply -f k8/deployment.yml && \
		kubectl apply -f k8/services.yml

deps:
	(docker network create chess || true ) && \
		docker run -d \
		--name chess_postgresql \
		--network chess \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=chess_api \
		-p 5432:5432 \
		postgres && \
	docker run -d \
		--name chess_nats \
		--network chess \
		-p 4222:4222 \
		nats

deps-clean:
	docker container stop chess_postgresql chess_nats || true && \
		docker container rm chess_postgresql chess_nats || true && \
		docker network rm chess || true

CHESS_API_DATABASE_USERNAME=postgres
CHESS_API_DATABASE_PASSWORD=password

CHESS_API_GITHUB_SECRET=secret
CHESS_API_GITHUB_KEY=key

CHESS_API_NATS_URL=localhost:4222
WORKDIR=${PWD}

run-local:
	export CHESS_API_DATABASE_USERNAME=$(CHESS_API_DATABASE_USERNAME) && \
	export CHESS_API_DATABASE_PASSWORD=$(CHESS_API_DATABASE_PASSWORD) && \
		export CHESS_API_GITHUB_SECRET=$(CHESS_API_GITHUB_SECRET) && \
		export CHESS_API_GITHUB_KEY=$(CHESS_API_GITHUB_KEY) && \
		CHESS_API_NATS_URL=$(CHESS_API_NATS_URL) && \
		go run $(WORKDIR)/cmd/chessapi/main.go start -c $(WORKDIR)/config/server.local.toml

run-chess-container-bash:
	docker run --rm -it \
		--network chess \
		-v $(WORKDIR):/go/src/github.com/dumbogo/chess \
		-p 8000:8000  \
		--name chessapi \
		--env-file .env \
		chess_chess bash


run-container-psql:
	docker exec -it chess_postgresql psql -U $(CHESS_API_DATABASE_USERNAME) chess_api
