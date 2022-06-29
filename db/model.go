package db

import "database/sql"

type Command struct {
	Id        int64
	Name      string
	Counter   sql.NullInt64
	Value     string
	IsCounter bool
	CreatedAt string
}

type FunFact struct {
	Id        int64
	CreatedAt string
	Value     string
}

type AnimalFact struct {
	Id        int64
	CreatedAt string
	Value     string
}

type Quote struct {
	Id        int64
	CreatedAt string
	Value     string
}
