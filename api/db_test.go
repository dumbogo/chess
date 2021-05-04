// +build integration

package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitDb(t *testing.T) { // TODO: initialize testdb programatically before running tests
	assert := assert.New(t)
	_, err := InitDb("localhost", "5432", "postgres", "password", "testdb")
	assert.Equal(nil, err)
}
