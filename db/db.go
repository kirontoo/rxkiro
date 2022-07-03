package db

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
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

func (db *Database) GetCommandByName(name string) Command {
	query := fmt.Sprintf(`SELECT * from "Commands" WHERE name = '%s' LIMIT 1`, name)
	row := db.Store.QueryRow(query)

	cmd := Command{}

	if err := row.Scan(&cmd.Id, &cmd.CreatedAt, &cmd.Name, &cmd.Counter, &cmd.Value, &cmd.IsCounter); err != nil {
		log.Logger.Error().Err(err).Send()
	}

	return cmd
}

func (db *Database) IncrementCounter(id int64, count int64) int64 {
	count = count + 1

	_, err := db.Store.Exec(fmt.Sprintf(`UPDATE "Commands" SET Counter = %d WHERE id = %d`, count, id))
	if err != nil {
		log.Logger.Error().Err(err).Send()
	}

	return count
}
