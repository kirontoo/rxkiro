package db

import (
	"database/sql"
	"errors"
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

func (db *Database) GetCommandByName(name string) (*Command, error) {
	query := fmt.Sprintf(`SELECT * from "Commands" WHERE name = '%s' LIMIT 1`, name)
	row := db.Store.QueryRow(query)

	cmd := Command{}

	if err := row.Scan(&cmd.Id, &cmd.CreatedAt, &cmd.Name, &cmd.Counter, &cmd.Value, &cmd.IsCounter); err != nil {
		log.Logger.Error().Err(err).Send()
		return nil, errors.New(fmt.Sprintf("Command, %s, does not exist in databse", name))
	}

	if cmd.Name != name {
		return nil, fmt.Errorf("Command, %s, does not exist in databse", name)
	}

	return &cmd, nil
}

func (db *Database) IncrementCounter(id int64, count int64) int64 {
	count = count + 1

	_, err := db.Store.Exec(
		fmt.Sprintf(`UPDATE "Commands" SET Counter = %d WHERE id = %d`, count, id),
	)
	if err != nil {
		log.Logger.Error().Err(err).Send()
	}
	return count
}

func (db *Database) AddCommand(name string, value string) string {
	exists := db.CommandExists(name)
	if !exists {
		return fmt.Sprintf("What are you doing!? Command '%s' already exists!", name)
	}

	query := fmt.Sprintf(`
	INSERT INTO "Commands" (
		name,
		value,
	)
	VALUES (
		%s,
		%s,
	)`, name, value)

	_, err := db.Store.Exec(query)
	if err != nil {
		log.Logger.Error().Err(err).Send()
		return "ERR: Could not create a new command. Contact the admin."
	}

	return fmt.Sprintf("New command added: %s", name)
}

func (db *Database) EditCommand(name string, value string) string {
	exists := db.CommandExists(name)
	if !exists {
		return fmt.Sprintf("Command '%s' doesn't exist", name)
	}

	query := fmt.Sprintf(`UPDATE "Commands" SET value = '%s' WHERE name = '%s'`, value, name)
	_, err := db.Store.Exec(query)
	if err != nil {
		log.Logger.Error().Err(err).Send()
		return "ERR Code 500. Amy go fix this."
	}

	return fmt.Sprintf("!%s was updated", name)
}

func (db *Database) DeleteCommand(name string) string {
	query := fmt.Sprintf(`DELETE FROM "Commands" WHERE name = '%s'`, name)
	_, err := db.Store.Exec(query)
	if err != nil {
		log.Logger.Error().Err(err).Send()
		return fmt.Sprintf("ERR: Could not delete command %s", name)
	}

	return fmt.Sprintf("!%s was deleted", name)
}

func (db *Database) CommandExists(name string) bool {
	cmd, _ := db.GetCommandByName(name)
	return cmd == nil
}

func (db *Database) AddCounter(name string) string {
	query := fmt.Sprintf(`
	INSERT INTO "Commands" (
		name,
		counter,
		value,
		isCounter,
	)
	VALUES (
		%s,
		0,
		"",
		true
	)`, name)

	_, err := db.Store.Exec(query)
	if err != nil {
		log.Logger.Error().Err(err).Send()
		return "ERR: Could not create a new counter. Contact the admin."
	}

	return fmt.Sprintf("New counter added: %s", name)
}
