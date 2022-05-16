package bot

import (
	"encoding/json"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/kirontoo/rxkiro/config"
	"github.com/kirontoo/rxkiro/db"
	"github.com/rs/zerolog"
	"github.com/supabase/postgrest-go"
)

type Command struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Counter   int64  `json:"counter"`
	Value     string `json:"value"`
	IsCounter bool   `json:"isCounter"`
	CreatedAt string `json:"created_at"`
}

type AnimalFact struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	Value     string `json:"value"`
}

type RxKiro struct {
	Client   *twitch.Client
	Config   config.Config
	Commands map[string]interface{}
	Log      zerolog.Logger
	db       postgrest.Client
}

const CommandTable = "Commands"

func NewBot(envPath string, logger zerolog.Logger) *RxKiro {
	botConfig, err := config.LoadConfig(envPath, "bot")
	if err != nil {
		logger.Error().Msgf("%s", err)
		logger.Fatal().Msg("Could not load env variables")
	}

	return &RxKiro{
		Commands: botCommands,
		Log:      logger,
		Config:   botConfig,
		Client:   twitch.NewClient(botConfig.BotName, botConfig.AuthToken),
		db:       db.Connect(botConfig.DbUrl, botConfig.DbToken),
	}
}

func trimQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func (b *RxKiro) Send(msg string) {
	b.Client.Say(b.Config.Streamer, msg)
}

func (b *RxKiro) RunCmd(cmdName string) {
	run, ok := b.Commands[cmdName]
	if ok {
		b.Log.Info().Str("cmd", cmdName).Msg("Executing Cmd")
		run.(func(*RxKiro))(b)
	} else {
		res, _, err := b.db.From(CommandTable).Select("*", "exact", false).Eq("name", cmdName).Execute()
		if err != nil {
			b.Log.Error().Msg(err.Error())
			return
		}

		var dbCmd []Command
		json.Unmarshal(res, &dbCmd)
		b.Log.Debug().Interface("Command", dbCmd).Msg("Cmds in DB")
		if len(dbCmd) > 0 {
			b.Log.Info().Str("cmd", cmdName).Msg("Cmd found in DB")
			b.Log.Info().Str("cmd", cmdName).Msg("Executing Cmd")

			data := dbCmd[0]

			if data.IsCounter {
				msg := b.runCounterCmd(data)
				b.Send(msg)
			} else if data.Value != "" {
				b.Send(data.Value)
			} else {
				b.Send("There's nothing here!")
			}
		} else {
			b.Log.Error().Str("cmd", cmdName).Msg("Invalid Cmd")
			b.Send("ERR: Invalid Command. Try again.")
		}
	}
}
