// +build integration

package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	dbHost     = "localhost"
	dbPort     = "5432"
	dbUser     = "postgres"
	dbPassword = "password"
	dbDatabase = "chess_api"
	dbSchema   = "public"
)

func dropSchema(gormDb *gorm.DB) { // TODO: implement
	r := gormDb.Exec(fmt.Sprintf("drop schema %s cascade", dbSchema))
	check(r.Error)
	r = gormDb.Exec(fmt.Sprintf("create schema %s", dbSchema))
	check(r.Error)
}

func truncateSchema(gormDb *gorm.DB) {
	r := gormDb.Exec(fmt.Sprintf("drop schema %s cascade", dbSchema))
	check(r.Error)
	r = gormDb.Exec(fmt.Sprintf("create schema %s", dbSchema))
	check(r.Error)
	Migrate(gormDb)
}

func initDbConnFactory() *gorm.DB {
	conn, err := InitDbConn(dbHost, dbPort, dbUser, dbPassword, dbDatabase)
	check(err)
	return conn
}

func TestInitDbConn(t *testing.T) {
	assert := assert.New(t)
	conn := initDbConnFactory()
	assert.NotNil(conn)
}

func TestMigrate(t *testing.T) {
	assert := assert.New(t)
	gormConn, err := InitDbConn(dbHost, dbPort, dbUser, dbPassword, dbDatabase)
	assert.Nil(err)
	dropSchema(gormConn)

	err = Migrate(gormConn)
	assert.Nil(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
