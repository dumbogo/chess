// +build integration

package api

import (
	context "context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func factoryServer() *Server {
	truncate()
	return &Server{
		Db: DBConn,
	}
}

func TestServerStartGame(t *testing.T) {
	assert := assert.New(t)
	server := factoryServer()
	ctxV, cancel := createCtxMetadataUser(&User{AccessToken: "hereistoken123", Email: "some@mail.com"})
	defer cancel()
	r, err := server.StartGame(ctxV, &StartGameRequest{
		Name:  "somename",
		Color: Color_WHITE,
	})
	assert.Nil(err)
	uuid.MustParse(r.GetUuid())
}

func TestServerJoinGame(t *testing.T) {
	assert := assert.New(t)
	server := factoryServer()
	ctx, cancel := createCtxMetadataUser(&User{AccessToken: "hereistoken123", Email: "some@mail.com"})
	defer cancel()
	// Create a game
	r, err := server.StartGame(ctx, &StartGameRequest{
		Name:  "somename",
		Color: Color_WHITE,
	})
	assert.Nil(err)
	uuid.MustParse(r.GetUuid())

	// Join a game:
	ctxJoin, cancelJoin := createCtxMetadataUser(&User{AccessToken: "someothertoken", Email: "other@mail.com"})
	defer cancelJoin()
	joinGameResponse, err := server.JoinGame(ctxJoin, &JoinGameRequest{
		Uuid: r.GetUuid(),
	})
	assert.Nil(err)
	assert.Equal(r.GetUuid(), joinGameResponse.GetUuid())
	assert.Equal(Color_BLACK, joinGameResponse.GetColor())
}

func TestServerMove(t *testing.T) {
	assert := assert.New(t)
	server := factoryServer()
	ctx, cancel := createCtxMetadataUser(&User{AccessToken: "hereistoken123", Email: "some@mail.com"})
	defer cancel()
	// Create a game
	r, err := server.StartGame(ctx, &StartGameRequest{
		Name:  "somename",
		Color: Color_WHITE,
	})
	assert.Nil(err)
	uuid.MustParse(r.GetUuid())

	// Join a game:
	ctxJoin, cancelJoin := createCtxMetadataUser(&User{AccessToken: "someothertoken", Email: "other@mail.com"})
	defer cancelJoin()
	joinGameResponse, err := server.JoinGame(ctxJoin, &JoinGameRequest{
		Uuid: r.GetUuid(),
	})
	assert.Nil(err)
	assert.Equal(r.GetUuid(), joinGameResponse.GetUuid())
	assert.Equal(Color_BLACK, joinGameResponse.GetColor())

	ctxMove, cancelMove := createCtxFromAccessToken("hereistoken123")
	defer cancelMove()
	moveResponse, err := server.Move(ctxMove, &MoveRequest{Uuid: r.GetUuid(), Color: Color_WHITE, FromSquare: "E2", ToSquare: "E4"})
	assert.Nil(err)
	assert.NotEmpty(moveResponse)

	ctxMove2, cancelMove2 := createCtxFromAccessToken("someothertoken")
	defer cancelMove2()
	moveResponse, err = server.Move(ctxMove2, &MoveRequest{Uuid: r.GetUuid(), Color: Color_BLACK, FromSquare: "G7", ToSquare: "G5"})
	assert.Nil(err)
	assert.NotEmpty(moveResponse)
}

func createCtxMetadataUser(u *User) (context.Context, context.CancelFunc) {
	tx := DBConn.Create(u)
	check(tx.Error)
	return createCtxFromAccessToken(u.AccessToken)
}

func createCtxFromAccessToken(t string) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	type keyctx string
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", t)})
	return metadata.NewIncomingContext(ctx, md), cancel
}
