package main

import (
	"fmt"
	"sync"

	"github.com/Adeithe/go-twitch/irc"
)

var chatHistory []string
var numLines int
var chatHistoryLock sync.Mutex
var maxLines int

func addToChat(msg irc.ChatMessage) {
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
