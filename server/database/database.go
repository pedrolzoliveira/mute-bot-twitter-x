package database

import (
	"flag"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var schema = `
CREATE TABLE IF NOT EXISTS muted_bots (
	user_handle VARCHAR(80) PRIMARY KEY
)
`

type MutedBots struct {
	UserHandle string `db:"user_handle"`
}

var db *sqlx.DB

func CreateDatabase() *sqlx.DB {
	var err error
	if db != nil {
		return db
	}

	if flag.Lookup("test.v") != nil {
		db, err = sqlx.Connect("sqlite3", ":memory:")
	} else {
		db, err = sqlx.Connect("sqlite3", "./mutebotx.db")
	}
	if err != nil {
		panic(err)
	}

	db.MustExec(schema)

	return db
}
