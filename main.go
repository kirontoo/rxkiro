package main

import (
	"log"
	"os"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/joho/godotenv"
)

var RxKiro = getEnvVariable("BOT_NAME")
var Streamer = getEnvVariable("STREAMER")
var OAUTH = getEnvVariable("OAUTH")

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	client := twitch.NewClient(RxKiro, OAUTH)

	client.Join(Streamer)
	defer client.Disconnect()

	log.Println("Bot has started")
	client.Say(Streamer, "/announce Hello World")

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Println(message.Message)
	})

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
