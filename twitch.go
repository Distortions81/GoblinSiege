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
		for x := 0; x < 10; x++ {
			startTime := time.Now()

			err := client.Connect()
			if err != nil {
				panic(err)
			}

			if time.Since(startTime) < time.Minute*10 {
				val := x * x * 2
				time.Sleep(time.Second * time.Duration(val))
			}
		}
	}()

	for x := 0; x < 3; x++ {
		msg := fmt.Sprintf("testing %v...\n", x)
		log.Println("Say: " + msg)
		client.Say(userSettings.Channel, msg)
		time.Sleep(time.Second * 5)
	}

}
