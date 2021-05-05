// +build integration

package api

import (
	context "context"
	"testing"
	"time"

	"github.com/dumbogo/chess/client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServerStartGame(t *testing.T) {
	assert := assert.New(t)
	conn, err := client.InitConn()
	assert.Nil(err)
	defer conn.Close()
	c := NewChessServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.StartGame(ctx, &StartGameRequest{
		Name:  "somename",
		Color: Color_WHITE,
	})
	assert.Nil(err)
	uuid.MustParse(r.GetUuid())
}
