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

func Connect(database string, connectStr string) *sql.DB {
	db, err := sql.Open(database, connectStr)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		log.Print("Db connected")
	}

	return db
}
