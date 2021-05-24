# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/dumbogo/chess
WORKDIR /go/src/github.com/dumbogo/chess

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go build -o chessapi /go/src/github.com/dumbogo/chess/cmd/api/main.go && \
    go build -o chess /go/src/github.com/dumbogo/chess/cmd/client/main.go
