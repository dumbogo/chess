// +build integration

package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "postgres"
	dbPassword = "password"
	dbDatabase = "chess_api"
	dbSchema   = "public"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	_, err := InitDbConn(dbHost, dbPort, dbUser, dbPassword, dbDatabase)
	check(err)
	truncate()
}

func truncate() {
	r := DBConn.Exec(fmt.Sprintf("drop schema %s cascade", dbSchema))
	check(r.Error)
	r = DBConn.Exec(fmt.Sprintf("create schema %s", dbSchema))
	check(r.Error)
	e := Migrate()
	check(e)
}

func TestInitDbConn(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(DBConn)
}

func TestMigrate(t *testing.T) {
	assert := assert.New(t)
	err := Migrate()
	assert.Nil(err)
}

func TestGetUserFromAccessToken(t *testing.T) {
	assert := assert.New(t)
	truncate()
	userCreated := User{
		AccessToken: "sometoken",
	}
	DBConn.Create(&userCreated)
	userFound := GetUserFromAccessToken("sometoken")
	assert.Equal(userCreated.ID, userFound.ID)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
