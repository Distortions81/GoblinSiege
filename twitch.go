package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch"
	"github.com/Adeithe/go-twitch/irc"
)

var Players map[int64]*playerData

type playerData struct {
	Points int
}

var chatHistory []string
var numLines int
var chatHistoryLock sync.Mutex

const maxLines = 20

func onShardMessage(shardID int, msg irc.ChatMessage) {
	//log.Printf("#%s %s: %s\n", msg.Channel, msg.Sender.ID, msg.Text)

	if !strings.EqualFold(msg.Channel, authInfo.Username) {
		//Ignore secondary channels
		log.Printf("Ignoring channel: %s -- %s\n", msg.Channel, authInfo.Username)
		return
	}

	if Players[msg.Sender.ID] == nil {
		log.Printf("Adding player '%v'\n", msg.Sender.ID)

		dbLock.Lock()
		Players[msg.Sender.ID] = &playerData{Points: 0}
		dbDirty = true
		dbLock.Unlock()
	}

	adminCommands(msg)

	//Add to chat history
	chatHistoryLock.Lock()
	defer chatHistoryLock.Unlock()
	out := fmt.Sprintf("%v: %v\n", msg.Sender.DisplayName, msg.Text)
	chatHistory = append(chatHistory, out)
	numLines++

	if numLines > maxLines {
		chatHistory = chatHistory[1:]
		numLines--
	}

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
