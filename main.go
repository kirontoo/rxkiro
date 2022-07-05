package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/kirontoo/rxkiro/bot"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const CMD_PREFIX = "!"

func main() {

	// Enable pretty logging
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("| %s |", i)
	}
	log.Logger = log.Output(output)

	rxkiro := bot.NewBot(".", log.Logger)

	// Connect to IRC
	rxkiro.Client.Join(rxkiro.Config.Streamer)

	// Listen for messages
	rxkiro.Client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Info().Str("Usr", message.User.DisplayName).Str("Msg", message.Message).Msg("New Message")
		cmd, ok := parseCommand(message.Message)
		if ok {
			log.Info().Str("cmd", cmd).Bool("ok", ok).Msg("Cmd Invoked")
			rxkiro.RunCmd(strings.ToLower(cmd), message)
		} else {
			log.Debug().Bool("ok?", ok).Send()
		}
	})

	rxkiro.Connect()

	// Disconnect from db and twitch chat
	defer rxkiro.Disconnect()
}

func parseCommand(message string) (string, bool) {
	messageWords := strings.SplitN(message, " ", 2)
	log.Print(messageWords)

	isCommand := strings.HasPrefix(messageWords[0], CMD_PREFIX)
	log.Debug().Str("Prefix", CMD_PREFIX).Bool("isCmd", isCommand).Msg("CommandPrefix")

	if isCommand {
		command := strings.TrimPrefix(messageWords[0], CMD_PREFIX)
		log.Debug().Str("Cmd", command).Msg("Command")
		return command, true
	}
	return "", false
}
