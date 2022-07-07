package bot

import (
	"fmt"
	"regexp"
	"strings"

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
	if cmdName == "cmd" {
		b.Send(b.handleCommandActions(message))
		return
	}

	if ok {
		// Execute hard coded command
		b.Log.Info().Str("cmd", cmdName).Msg("Executing Cmd")
		run.(func(*RxKiro))(b)
	} else {
		cmd, _ := b.db.GetCommandByName(cmdName)
		if cmd != nil {
			if !cmd.IsCounter {
				// check if any cmd vars need to be replaced
				matches := b.findCmdVars(cmd.Value)
				if len(matches) > 0 {
					updatedMsg := cmd.Value
					for _, m := range matches {
						updatedMsg = b.replaceCmdVariables(m, updatedMsg, message)
					}

					b.Send(updatedMsg)
				} else {
					b.Send(cmd.Value)
				}
			} else {
				msg := b.updateCounter(cmd)
				b.Send(msg)
			}
		} else {
			b.Send("ERR: Invalid Command. Try again.")
			b.Log.Error().Str("cmd", cmdName).Msg("Invalid Cmd")
		}
	}
}

func (b *RxKiro) handleCommandActions(msg twitch.PrivateMessage) string {
	// ["!cmd", "add", "cmdname", "some", "words", "here"]
	raw := strings.Split(msg.Message, " ")
	params := len(raw)

	errorMsg := "ERR: Invalid syntax. Try !cmd add test this is a test command LUL"
	if params < 3 {
		return errorMsg
	}

	action := raw[1]
	name := strings.ToLower(raw[2])
	rawMessage := strings.Join(raw[3:], " ")

	if (name == "add" || name == "edit") && params < 4 {
		return errorMsg
	}

	b.Log.Debug().Str("cmd action", action).Str("name", name).Str("value", rawMessage).Msg("")

	switch action {
	case "add":
		output := b.db.AddCommand(name, rawMessage)
		b.Send(output)
	case "edit":
		return b.db.EditCommand(name, rawMessage)
	case "delete":
		if params < 3 {
			return "ERR: Invalid syntax. Try !cmd delete test"
		}
		return b.db.DeleteCommand(name)
	default:
		return errorMsg
	}

	// TODO:
	return ""
}

func (b *RxKiro) updateCounter(cmd *db.Command) string {
	if cmd.Counter.Valid {
		count := b.db.IncrementCounter(cmd.Id, cmd.Counter.Int64)
		return fmt.Sprintf("%s: %d", cmd.Name, count)
	}

	return fmt.Sprintf("Invalid counter: %s", cmd.Name)
}

func (b *RxKiro) replaceCmdVariables(rawCmd string, s string, msg twitch.PrivateMessage) string {
	raw := strings.ToLower(strings.Trim(rawCmd, "${}"))
	cmdVars := strings.Split(raw, " ")
	cmdVar := cmdVars[0]

	b.Log.Debug().Str("cmdvar", cmdVar).Send()
	b.Log.Print(cmdVars)

	switch cmdVar {
	case "user":
		username := "@" + msg.User.DisplayName
		updatedMsg := strings.ReplaceAll(s, rawCmd, username)
		return updatedMsg
	case "mention":
		var updatedMsg string
		m := strings.Split(msg.Message, " ")
		if len(m) < 2 {
			updatedMsg = strings.ReplaceAll(s, rawCmd, "")
		} else {
			mentionedUser := m[1]
			updatedMsg = strings.ReplaceAll(s, rawCmd, mentionedUser)
		}
		return updatedMsg
	case "random":
		return randomCmd(s, rawCmd)
	default:
		return ""
	}
}

func (b *RxKiro) findCmdVars(s string) []string {
	var matches = []string{}

	r, err := regexp.Compile(`\${(.*?)}`)
	if err != nil {
		b.Log.Error().Err(err).Send()
	}

	matched := r.Match([]byte(s))
	if matched {
		matches = r.FindAllString(s, -1)
	}

	return matches
}
