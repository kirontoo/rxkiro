package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/kirontoo/rxkiro/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/supabase/postgrest-go"
)

var botConfig config.Config
var dbClient postgrest.Client
var commands []Command

// var bot *twitch.Client

const CMD_PREFIX = "!"

type Command struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Counter   int64  `json:"counter"`
	Value     string `json:"value"`
	IsCount   bool   `json:"isCount"`
	CreatedAt string `json:"created_at"`
}

type RxKiro struct {
	client *twitch.Client
	config config.Config
}

func main() {
	var bot RxKiro

	// Enable pretty logging
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("| %s |", i)
	}
	log.Logger = log.Output(output)

	botConfig, err := config.LoadConfig(".", "bot")
	if err != nil {
		log.Error().Msgf("%s", err)
		log.Fatal().Msg("Could not load env variables")
	}

	bot.config = botConfig

	// dbClient := db.Connect(botConfig.DbUrl, botConfig.DbToken)
	// data, _, rqerr := dbClient.From("Commands").Select("*", "", false).Execute()
	// if rqerr != nil {
	// 	log.Error().Msg(rqerr.Error())
	// }

	// json.Unmarshal(data, &commands)
	// // log.Print(commands, countType)
	// log.Debug().Interface("commands", commands).Send()

	bot.client = twitch.NewClient(bot.config.BotName, bot.config.AuthToken)

	// Connect to IRC
	bot.client.Join(botConfig.Streamer)

	log.Info().Str("Streamer", botConfig.Streamer).Str("Bot name", botConfig.BotName).Msg("Connected to chat")
	bot.send("/announce Hello World")

	// Listen for messages
	bot.client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		cmd, ok := parseCommand(message.Message)
		if ok {
			log.Info().Str("cmd", cmd).Bool("ok", ok).Msg("Found Command")
			bot.runCmd(cmd)
		}
		return
	})

	clientErr := bot.client.Connect()
	defer bot.client.Disconnect()
	if clientErr != nil {
		panic(clientErr)
	}
}

func parseCommand(message string) (string, bool) {
	messageWords := strings.SplitN(message, " ", 2)

	isCommand := strings.HasPrefix(messageWords[0], CMD_PREFIX)
	log.Debug().Str("Prefix", CMD_PREFIX).Msg("CommandPrefix")

	if isCommand {
		command := strings.TrimPrefix(messageWords[0], CMD_PREFIX)
		log.Debug().Str("Cmd", command)
		return command, true
	}
	return "", false
}

func (b *RxKiro) send(msg string) {
	b.client.Say(b.config.Streamer, msg)
}

func (b *RxKiro) runCmd(cmd string) {
	switch cmd {
	case "ping":
		b.send("pong")
		b.send("where is this going?")
	default:
		b.send("do you....need help?")
	}
}
