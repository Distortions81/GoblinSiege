package main

import (
	"fmt"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

var client *twitch.Client

func fSay(format string, args ...interface{}) {
	buf := fmt.Sprintf(format, args...)
	client.Say(userSettings.Channel, buf)
}

func connectTwitch() {
	readSettings()

	client = twitch.NewClient(userSettings.UserName, "oauth:"+userSettings.AuthToken)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		qlog("%v: %v\n", message.User.DisplayName, message.Message)
	})

	qlog("Joining channel: %v", userSettings.Channel)
	client.Join(userSettings.Channel)

	qlog("Connecting to twitch...")
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
}
