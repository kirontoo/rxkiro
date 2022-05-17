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


## License

This project is licensed under the [MIT](LICENSE.md) License - see the [LICENSE.md](LICENSE.md) file for
details
