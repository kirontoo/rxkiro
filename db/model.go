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

// TODO: FunFact and AnimalFact can be merged to one struct called Fact

type Fact struct {
	Id        int64
	CreatedAt string
	Value     string
}

type Quote struct {
	Id        int64
	CreatedAt string
	Value     string
}
