package bot

import (
	"database/sql"
	"fmt"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/kirontoo/rxkiro/config"
	"github.com/kirontoo/rxkiro/db"
	"github.com/rs/zerolog"
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
	db       *sql.DB
}

const CommandTable = "Commands"

func NewBot(envPath string, logger zerolog.Logger) *RxKiro {
	botConfig, err := config.LoadConfig(envPath, "bot")
	if err != nil {
		logger.Error().Msgf("%s", err)
		logger.Fatal().Msg("Could not load env variables")
	}

	// connStr := "user=postgres password=x2yLkNsMA5ZQAig9 host=db.ddtrcmquqodwxhmkvflt.supabase.co port=5432 dbname=postgres"
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		botConfig.DbHost, botConfig.DbPort, botConfig.DbUser, botConfig.DbPassword, botConfig.DbName)

	return &RxKiro{
		Log:      logger,
		Config:   botConfig,
		Client:   twitch.NewClient(botConfig.BotName, botConfig.AuthToken),
		db:       db.Connect("postgres", connStr),
		Commands: botCommands,
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

func (b *RxKiro) Disconnect() {
	b.Client.Disconnect()
	b.db.Close()
}

func (b *RxKiro) Send(msg string) {
	b.Client.Say(b.Config.Streamer, msg)
}

func (b *RxKiro) RunCmd(cmdName string, message twitch.PrivateMessage) {
	run, ok := b.Commands[cmdName]
	b.Log.Debug().Str("cmd", cmdName).Bool("ok?", ok).Msg("RunCmd")

	if ok {
		// Execute hard coded command
		b.Log.Info().Str("cmd", cmdName).Msg("Executing Cmd")
		run.(func(*RxKiro))(b)
	} else {
		// // Find & Execute commands from the database
		// res, _, err := b.db.From(CommandTable).Select("*", "exact", false).Eq("name", cmdName).Execute()
		// if err != nil {
		// 	b.Log.Error().Msg(err.Error())
		// 	return
		// }

		// var dbCmd []Command
		// json.Unmarshal(res, &dbCmd)
		// b.Log.Debug().Interface("Command", dbCmd).Msg("Cmds in DB")
		// if len(dbCmd) > 0 {
		// 	b.Log.Info().Str("cmd", cmdName).Msg("Cmd found in DB")
		// 	b.Log.Info().Str("cmd", cmdName).Msg("Executing Cmd")

		// 	data := dbCmd[0]

		// 	if data.IsCounter {
		// 		msg := b.runCounterCmd(data)
		// 		b.Send(msg)
		// 	} else if data.Value != "" {
		// 		// check if any vars need to be replaced
		// 		matches := b.findCmdVars(data.Value)
		// 		if len(matches) > 0 {
		// 			for _, m := range matches {
		// 				// cmdVar := strings.ToLower(strings.Trim(m, "${}"))
		// 				newMsg := b.replaceCmdVariables(m, data.Value, message)
		// 				if newMsg != "" {
		// 					b.Send(newMsg)
		// 				}
		// 			}
		// 		} else {
		// 			b.Send(data.Value)
		// 		}
		// 	} else {
		// 		b.Send("There's nothing here!")
		// 	}
		// } else {
		// 	b.Log.Error().Str("cmd", cmdName).Msg("Invalid Cmd")
		// 	// b.Send("ERR: Invalid Command. Try again.")
		// }
	}
}
