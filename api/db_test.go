// +build integration

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDbConn(t *testing.T) {
	assert := assert.New(t)
	_, err := InitDbConn("localhost", "5432", "postgres", "password", "chess_api")
	assert.Equal(nil, err)
}
