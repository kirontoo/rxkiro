package bot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/kirontoo/rxkiro/config"
	"github.com/kirontoo/rxkiro/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type RxKiro struct {
	Client   *twitch.Client
	Config   config.Config
	Commands map[string]interface{}
	Log      zerolog.Logger
	db       *db.Database
}

const CommandTable = "Commands"

func NewBot(envPath string, logger zerolog.Logger) *RxKiro {
	botConfig, err := config.LoadConfig(envPath, "bot")
	if err != nil {
		logger.Error().Msgf("%s", err)
		logger.Fatal().Msg("Could not load env variables")
	}

	return &RxKiro{
		Log:      logger,
		Config:   botConfig,
		Client:   twitch.NewClient(botConfig.BotName, botConfig.AuthToken),
		db:       nil,
		Commands: botCommands,
	}
}

func (b *RxKiro) Connect() {

	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		b.Config.DbHost, b.Config.DbPort, b.Config.DbUser, b.Config.DbPassword, b.Config.DbName)
	b.db = db.Connect("postgres", connStr)

	twtErr := b.Client.Connect()
	if twtErr != nil {
		b.Log.Panic().Msg(twtErr.Error())
	} else {
		log.Info().Str("Streamer", b.Config.Streamer).Str("Bot name", b.Config.BotName).Msg("Connected to chat")
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

	if ok {
		// Execute hard coded command
		b.Log.Info().Str("cmd", cmdName).Msg("Executing Cmd")
		run.(func(*RxKiro))(b)
	} else {
		cmd := b.db.GetCommandByName(cmdName)
		if cmd.Name == cmdName {
			if !cmd.IsCounter {
				// check if any cmd vars need to be replaced
				matches := b.findCmdVars(cmd.Value)
				if len(matches) > 0 {
					for _, m := range matches {
						newMsg := b.replaceCmdVariables(m, cmd.Value, message)
						if newMsg != "" {
							b.Send(newMsg)
						}
					}
				} else {
					b.Send(cmd.Value)
				}
			} else {
				msg := b.runCounter(cmd)
				b.Send(msg)
			}
		} else {
			b.Send("ERR: Invalid Command. Try again.")
			b.Log.Error().Str("cmd", cmdName).Msg("Invalid Cmd")
		}
	}
}
