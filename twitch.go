package main

import (
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

func connectTwitch() {
	readSettings()

	log.Println("Connecting to twitch...")

	client := twitch.NewClient(userSettings.UserName, "oauth:"+userSettings.AuthToken)

	err := client.Connect()
	if err != nil {
		panic(err)
	}

	client.Join(userSettings.Channel)

	for x := 0; x < 10; x++ {
		msg := "testing 123..."
		log.Printf("Say: %v\n", msg)
		client.Say(userSettings.Channel, msg)
		time.Sleep(time.Second * 5)
	}
}
