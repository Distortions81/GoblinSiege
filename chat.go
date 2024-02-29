package main

import (
	"fmt"
	"sync"

	"github.com/Adeithe/go-twitch/irc"
)

var chatHistory []string
var numLines int
var chatHistoryLock sync.Mutex
var maxShowLines int

const maxHistory = 100

func addToChat(msg irc.ChatMessage) {
	chatHistoryLock.Lock()
	defer chatHistoryLock.Unlock()

	out := fmt.Sprintf("%v: %v\n", msg.Sender.DisplayName, msg.Text)
	chatHistory = append(chatHistory, out)
	numLines++

	trimChatHistory()
}

func trimChatHistory() {
	//Remove old lines until we fit
	for {
		if numLines > maxHistory {
			chatHistory = chatHistory[1:]
			numLines--
		} else {
			break
		}
	}
}
