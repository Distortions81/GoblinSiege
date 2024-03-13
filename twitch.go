package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

const maxConsecutiveReconnects = 10
const reconnectLimiterTimeout = time.Minute * 10

var client *twitch.Client

func fSay(format string, args ...interface{}) {
	buf := fmt.Sprintf(format, args...)
	client.Say(userSettings.UserName, buf)
}

func connectTwitch() {

	client = twitch.NewClient(userSettings.UserName, "oauth:"+userSettings.AuthToken)

	client.OnPrivateMessage(handleChat)

	qlog("Joining channel: %v", userSettings.UserName)
	client.Join(userSettings.UserName)

	qlog("Connecting to twitch...")
	go func() {
		for x := 0; x < maxConsecutiveReconnects; x++ {
			startTime := time.Now()

			err := client.Connect()
			if err != nil {
				panic(err)
			}

			if time.Since(startTime) < reconnectLimiterTimeout {
				val := x * x * 2 //Increasing delay
				time.Sleep(time.Second * time.Duration(val))
			} else {
				//We were connected long enough that we can reset the reconnect limiter
				x = 0
			}
		}
	}()
}

func strToID(input string) int64 {
	userid, _ := strconv.ParseInt(input, 10, 64)
	return userid
}
