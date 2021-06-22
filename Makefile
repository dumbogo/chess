# Go command to use for build
GO ?= go

BINDIR ?= bin
RELEASEDIR ?= releases

## HARDCODED, needs to automatically detects current version by git
version ?= v1.0.0-alpha.1

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
	go build -o $(BINDIR)/chess cmd/chess/main.go
	GOOS=linux go build -o $(BINDIR)/chessapi cmd/chessapi/main.go

clean:
	rm -rf bin/* && rm -rf releases/*.tar.gz releases/sha256sums.txt releases/*.zip

proto: # Build proto files
	GO111MODULE=on protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/service.proto
