package db

import (
	"database/sql"
	"log"
)

type Database struct {
	Store *sql.DB
}

func Connect(database string, connectStr string) *Database {
	db, err := sql.Open(database, connectStr)
	if err != nil {
		log.Print("db connect: panic")
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		log.Print("Db connected")
	}

	return &Database{
		Store: db,
	}
}

func (db *Database) Close() {
	db.Close()
}
