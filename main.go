package main

import (
	"log"

	"github.com/gempir/go-twitch-irc/v3"
)

const RxKiro = "rxKiro"
const Streamer = "kironto"

func main() {
	// TODO: make this an env value
	client := twitch.NewClient(RxKiro, "oauth:")

	client.Join(Streamer)
	defer client.Disconnect()

	log.Println("Bot has started")
	client.Say(Streamer, "Hello World")

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		log.Println(message.Message)
	})

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}
