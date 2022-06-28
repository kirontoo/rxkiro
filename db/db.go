package db

import (
	"database/sql"
	"log"
)

// func Connect(url string, token string) postgrest.Client {
// 	dbClient := postgrest.NewClient(url, "", nil)
// 	if dbClient.ClientError != nil {
// 		panic(dbClient.ClientError)
// 	}

// 	dbClient = dbClient.TokenAuth(token)

// 	return *dbClient
// }

type Database struct {
	SqlDb *sql.DB
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
		SqlDb: db,
	}
}

func (db *Database) Close() {
	db.Close()
}
