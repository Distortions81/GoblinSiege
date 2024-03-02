package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

var client *twitch.Client

func connectTwitch() {
	readSettings()

	client = twitch.NewClient(userSettings.UserName, "oauth:"+userSettings.AuthToken)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Printf("%v: %v\n", message.User.DisplayName, message.Message)
	})

	log.Printf("Joining channel: %v\n", userSettings.Channel)
	client.Join(userSettings.Channel)

	log.Println("Connecting to twitch...")
	go func() {
		err := client.Connect()
		if err != nil {
			panic(err)
		}
	}()

	for x := 0; x < 10; x++ {
		msg := "testing 123..."
		log.Printf("Say: %v\n", msg)
		client.Say(userSettings.Channel, msg)
		time.Sleep(time.Second * 5)
	}

}
