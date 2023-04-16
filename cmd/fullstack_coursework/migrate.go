package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const (
	fixturesDir = "./fixtures"
	driver      = "postgres"
)

func migrate(postgresUrl, command string) error {
	db, err := sql.Open(driver, postgresUrl)
	if err != nil {
		return err
	}
	return goose.Run(command, db, fixturesDir)
}
