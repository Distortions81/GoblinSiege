package main

import (
	"log"
	"strings"
	"time"

	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/irc"
)

var twitchWriter = &irc.Conn{}
var twitchReader = &irc.Client{}

func onShardMessage(shardID int, msg irc.ChatMessage) {

	if !strings.EqualFold(msg.Channel, userSettings.UserName) {
		//Ignore secondary channels
		return
	}

	log.Printf("#%v: %v: %v\n", msg.Channel, msg.Sender.DisplayName, msg.Text)

	handleChat(msg)
}

func connectTwitch() {
	readSettings()

	//Connect
	twitchWriter.SetLogin(userSettings.UserName, "oauth:"+string(userSettings.AuthToken))
	if err := twitchWriter.Connect(); err != nil {
		panic("failed to start writer")
	}

	twitchReader = twitch.IRC()
	twitchReader.OnShardReconnect(onShardReconnect)
	twitchReader.OnShardLatencyUpdate(onShardLatencyUpdate)
	twitchReader.OnShardMessage(onShardMessage)

	if err := twitchReader.Join(userSettings.UserName); err != nil {
		panic(err)
	}

	log.Println("Connected to IRC!")
}

func onShardReconnect(shardID int) {
	log.Printf("Shard #%d reconnected\n", shardID)
}

func onShardLatencyUpdate(shardID int, latency time.Duration) {
	log.Printf("Shard #%d has %dms ping\n", shardID, latency.Milliseconds())
}
