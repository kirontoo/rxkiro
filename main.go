package main

import (
	"fmt"
	"os"
	"time"

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
	rxkiro.Join()
	rxkiro.ListenToChat()
	rxkiro.Connect()

	// Disconnect from db and twitch chat
	defer rxkiro.Disconnect()
}
