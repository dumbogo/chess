// +build integration

package api

import (
	context "context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func factoryServer() *Server {
	dbConn := initDbConnFactory()
	return &Server{
		Db: dbConn,
	}
}

func TestServerStartGame(t *testing.T) {
	assert := assert.New(t)
	server := factoryServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := server.StartGame(ctx, &StartGameRequest{
		Name:  "somename",
		Color: Color_WHITE,
	})
	assert.Nil(err)
	uuid.MustParse(r.GetUuid())
}

func TestServerJoinGame(t *testing.T) {
	// Create a game
	assert := assert.New(t)
	server := factoryServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := server.StartGame(ctx, &StartGameRequest{
		Name:  "somename",
		Color: Color_WHITE,
	})
	assert.Nil(err)
	uuid.MustParse(r.GetUuid())

	// Join a game:
	ctxJoin, cancelJoin := context.WithTimeout(context.Background(), time.Second)
	defer cancelJoin()
	joinGameResponse, err := server.JoinGame(ctxJoin, &JoinGameRequest{
		Uuid: r.GetUuid(),
	})
	assert.Nil(err)
	assert.Equal(r.GetUuid(), joinGameResponse.GetUuid())
	assert.Equal(Color_BLACK, joinGameResponse.GetColor())
}
