package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var RxKiro = getEnvVariable("BOT_NAME")
var Streamer = getEnvVariable("STREAMER")
var OAUTH = getEnvVariable("OAUTH")

const COMMAND_PREFIX = "!"

func main() {
	// Enable pretty logging
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC822}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("| %s |", i)
	}

	log.Logger = log.Output(output)

	// Create a new twitch IRC client
	client := twitch.NewClient(RxKiro, OAUTH)

	// Connect to IRC
	client.Join(Streamer)
	defer client.Disconnect()

	log.Info().Str("Streamer", Streamer).Str("Bot name", RxKiro).Msg("Connected to chat")
	client.Say(Streamer, "/announce Hello World")

	// Listen for messages
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Printf(message.Message)
		parseCommand(message)
	})

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	return os.Getenv(key)
}

func parseCommand(message twitch.PrivateMessage) {
	messageWords := strings.Split(message.Message, " ")

	isCommand := strings.HasPrefix(messageWords[0], COMMAND_PREFIX)

	if isCommand {
		command := strings.TrimPrefix(messageWords[0], COMMAND_PREFIX)
		log.Printf("Command: %s", command)
	}

	return
}
