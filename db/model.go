package db

type Command struct {
	id        int64
	Name      string
	Counter   int64
	Value     string
	IsCounter bool
	CreatedAt string
}

type FunFact struct {
	id        int64
	CreatedAt string
	Value     string
}

type AnimalFact struct {
	id        int64
	CreatedAt string
	Value     string
}

type Quote struct {
	id        int64
	CreatedAt string
	Value     string
}
