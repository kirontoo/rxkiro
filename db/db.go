package db

import "github.com/supabase/postgrest-go"

func Connect(url string, token string) postgrest.Client {
	dbClient := postgrest.NewClient(url, "", nil)
	if dbClient.ClientError != nil {
		panic(dbClient.ClientError)
	}

	dbClient = dbClient.TokenAuth(token)

	return *dbClient
}
