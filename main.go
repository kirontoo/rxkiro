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
)

var botConfig config.Config

func main() {
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

	// Create a new twitch IRC client
	client := twitch.NewClient(botConfig.BotName, botConfig.AuthToken)

	// Connect to IRC
	client.Join(botConfig.Streamer)
	defer client.Disconnect()

	log.Info().Str("Streamer", botConfig.Streamer).Str("Bot name", botConfig.BotName).Msg("Connected to chat")
	client.Say(botConfig.Streamer, "/announce Hello World")

	// Listen for messages
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Printf(message.Message)
		parseCommand(message)
	})

	clientErr := client.Connect()
	if clientErr != nil {
		panic(clientErr)
	}
}

func parseCommand(message twitch.PrivateMessage) {
	messageWords := strings.Split(message.Message, " ")

	isCommand := strings.HasPrefix(messageWords[0], botConfig.CmdPrefix)

	if isCommand {
		command := strings.TrimPrefix(messageWords[0], botConfig.CmdPrefix)
		log.Debug().Str("Cmd", command)
	}

	return
}
