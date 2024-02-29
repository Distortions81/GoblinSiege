package main

import (
	"log"
	"strings"
	"time"

	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/irc"
)

func onShardMessage(shardID int, msg irc.ChatMessage) {

	if !strings.EqualFold(msg.Channel, userSettings.Username) {
		//Ignore secondary channels
		return
	}

	log.Printf("%s: %s\n", msg.Sender.DisplayName, msg.Text)

	if Players[msg.Sender.ID] == nil {
		log.Printf("Adding player '%v' to db.\n", msg.Sender.ID)

		dbLock.Lock()
		Players[msg.Sender.ID] = &playerData{Points: 0}
		dbDirty = true
		dbLock.Unlock()
	}

	if adminCommands(msg) {
		return
	}

	handleChat(msg)

}

func connectTwitch() {
	writer := &irc.Conn{}

	readSettings()

	//Connect
	writer.SetLogin(userSettings.Username, "oauth:"+string(userSettings.AuthToken))
	if err := writer.Connect(); err != nil {
		panic("failed to start writer")
	}

	reader := twitch.IRC()
	reader.OnShardReconnect(onShardReconnect)
	reader.OnShardLatencyUpdate(onShardLatencyUpdate)
	reader.OnShardMessage(onShardMessage)

	if err := reader.Join(userSettings.Username); err != nil {
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
