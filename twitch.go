package main

import (
	"log"
	"strings"
	"time"

	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/irc"
)

var Players map[string]*playerData

type playerData struct {
	Points int
}

func onShardMessage(shardID int, msg irc.ChatMessage) {
	//log.Printf("#%s %s: %s\n", msg.Channel, msg.Sender.Username, msg.Text)

	if !strings.EqualFold(msg.Channel, authInfo.Username) {
		//Ignore secondary channels
		log.Printf("Ignoring channel: %s -- %s\n", msg.Channel, authInfo.Username)
		return
	}

	if len(msg.Sender.Username) > 0 && Players[msg.Sender.Username] == nil {
		log.Printf("Adding player '%v'\n", msg.Sender.Username)

		dbLock.Lock()
		Players[msg.Sender.Username] = &playerData{Points: 0}
		dbDirty = true
		dbLock.Unlock()
	}

	log.Printf("%s: %s\n", msg.Sender.Username, msg.Text)
}

func connectTwitch() {
	writer := &irc.Conn{}

	readAuth()

	//Connect
	writer.SetLogin(authInfo.Username, "oauth:"+string(authInfo.AuthToken))
	if err := writer.Connect(); err != nil {
		panic("failed to start writer")
	}

	reader := twitch.IRC()
	reader.OnShardReconnect(onShardReconnect)
	reader.OnShardLatencyUpdate(onShardLatencyUpdate)
	reader.OnShardMessage(onShardMessage)

	if err := reader.Join(authInfo.Username); err != nil {
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
