package main

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/Adeithe/go-twitch/irc"
)

type chatMessageData struct {
	message string
	color   color.Color
	len     int
}

var chatHistory []chatMessageData
var numLines int
var chatHistoryLock sync.Mutex
var maxShowLines int

const maxHistory = 100

func addToChat(msg irc.ChatMessage) {
	chatHistoryLock.Lock()
	defer chatHistoryLock.Unlock()

	out := fmt.Sprintf("%v: %v\n", msg.Sender.DisplayName, msg.Text)
	c, err := Hex2Color(msg.Sender.Color)

	var msgColor = color.RGBA{R: 255, G: 255, B: 255}
	if err == nil {
		msgColor = c
	}

	chatHistory = append(chatHistory, chatMessageData{message: out, color: msgColor, len: len(out)})
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
