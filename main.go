package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/kirontoo/rxkiro/bot"
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

	log.Info().Str("Streamer", rxkiro.Config.Streamer).Str("Bot name", rxkiro.Config.BotName).Msg("Connected to chat")
	rxkiro.Send("/announce Hello World")

	// Listen for messages
	rxkiro.Client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		cmd, ok := parseCommand(message.Message)
		if ok {
			log.Info().Str("cmd", cmd).Bool("ok", ok).Msg("New Cmd Found")
			rxkiro.RunCmd(cmd)
		}
		return
	})

	ClientErr := rxkiro.Client.Connect()
	defer rxkiro.Client.Disconnect()
	if ClientErr != nil {
		panic(ClientErr)
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
