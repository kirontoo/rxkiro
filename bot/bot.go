package bot

import (
	"encoding/json"
	"strconv"

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
		Commands: commands,
		Log:      logger,
		Config:   botConfig,
		Client:   twitch.NewClient(botConfig.BotName, botConfig.AuthToken),
		db:       db.Connect(botConfig.DbUrl, botConfig.DbToken),
	}
}

var commands = map[string]interface{}{
	"ping": func(b *RxKiro) {
		// TODO after pinged 5 times in a row, send back a special msg
		// might need to save a the time and a db record for this
		b.Send("pong")
		b.Send("where is this going?")
	},
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
				data.Counter = b.incrementCounter(data.Counter)
				b.Log.Debug().Int64("counter", data.Counter).Msg("Should increment")
				updatedRes, _, err := b.db.From(CommandTable).Update(data, "representation", "exact").Eq("name", data.Name).Execute()

				var jsonRes Command
				json.Unmarshal(updatedRes, &jsonRes)
				if err != nil {
					b.Log.Error().Interface("Command", jsonRes).Msg("Error updating")
					b.Log.Error().Err(err)
				}

				msg := data.Name + ": " + strconv.Itoa(int(data.Counter))
				b.Send(msg)
			}

			if data.Value != "" {
				b.Send(data.Value)
			} else {
				b.Send("There's nothing here!")
			}
		} else {
			b.Log.Error().Str("cmd", cmdName).Msg("Invalid Cmd")
			b.Send("This is invalid!")
		}
	}
}

func (b *RxKiro) incrementCounter(count int64) int64 {
	return (count + 1)
}
