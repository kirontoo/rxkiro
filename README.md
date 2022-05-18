# RxKiro

A simple twitch bot for my stream.


## Getting Started

### Prerequisites


- [Go](https://go.dev/doc/install)
- [Supabase](https://supabase.com/) 
- [Twitch Account](https://twitch.tv)


### Installing

Install [Go](https://go.dev/doc/install) on your machine and then setup a [Supabase](https://supabase.com/) project for the bot.
This bot uses a PostgreSQL database as its backend. 

**Step 1**: Clone the repo
```bash
// with https
git clone https://github.com/kirontoo/rxkiro.git

// ssh
git clone git@github.com:kirontoo/rxkiro.git
```

**Step 2**: Make sure to install all Go modules
```
go mod download
```

### Setting Up

#### Environment Variables
Make a environment file called `bot.env`.

Here are all the environment variables you need to get started.
```
AUTH_TOKEN=your twitch auth login token
BOT_NAME=bot name or user you are logging in as
STREAMER=channel name (streamer you want to connet the bot to)
DB_API_URL=supabase rest api url
DB_TOKEN=supabase api token
```

- `AUTH_TOKEN` is the twitch auth token. You can grab one with this [link](https://twitchapps)
- `BOT_NAME` is the username of your twitch account that you want to log in as
- `STREAMER` is the channel name/streamer you want to connect the bot to
- `DB_API_URL` is the supabase rest api url
- `DB_TOKEN` is the supabase api token

To find the `DB_API_URL` and the `DB_TOKEN`, you need to go into the `Settings` of your Supabase project.
- Under `Project API Keys`, copy the `anon public` key and set it as the `DB_TOKEN`. 
- Below this in `Configuration`, copy the `URL` and set it as the `DB_API_URL`

Here is a official [guide](https://supabase.com/docs/guides/api#api-url-and-keys) if my instructions are confusing.

#### Setting Up Your Database

The bot currently only uses the `Commands` table to look for commands. The bot 
will first check for default commands that's been hard coded. If 
that command is not found, then it will check the database for the correct command.


The `Commands` table should have these columns:
`id, name, counter, isCounter, value, created_at`
```go
type Command struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    Counter   int64  `json:"counter"`
    Value     string `json:"value"`
    IsCounter bool   `json:"isCounter"`
    CreatedAt string `json:"created_at"`
}
```

## Running the Bot
Run this comand

```go
go run main.go
```

## Built With
- [Go](https://go.dev/doc/install)
- [Supabase](https://supabase.com/) - Used for serverless SQL database

## Features

### Command variables
Right now, there is only 1 valid command variable `user` which will replace 
command variable with the username of the user who invoked the command.

**For Example**: 
`!lurk` has the value of `${user} is lurking!`. The `!lurk` command will 
replace `${user}` with a username. The output will look like this: `Kironto is lurking!`.

### Random Facts
There is a built in command for `animal facts` and `fun facts`, they are the 
same but uses different tables depending on which type of fact you want.
In the future I might combine these tables and add a column for 'fact type'.

**Here are columns that need to be defined for these tables**:
```json
{
    id: int8,
    created_at: timestampz,
    value: text
}
```

These commands can be invoked in the chat with `!animalfact` or `!mefact`.
For now, to properly use these commands, you'll need to create a custom function in `Supabase`.

In the `SQL Editor` tab of the Supabase dashboard, create a new query called "Rand Animal Fact"
and then paste this code.

```sql
create or replace function rand_animal_fact()
returns text
language sql
as $$
  SELECT value FROM "AnimalFacts" ORDER BY RANDOM() limit 1;
$$
```

Run this query in the SQL editor and the custom function should be added to 
your db. Do the same for the fun facts except replace `rand_animal_fact` with 
`rand_fun_fact` and `AnimalFacts` with `FunFact`.

`!animalfact` and `!mefact` commands should work now.

## License

This project is licensed under the [MIT](LICENSE.md) License - see the [LICENSE.md](LICENSE.md) file for
details
