package main

import (
	"goTwitchGame/sclean"
	"image/color"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch/irc"
)

type chatMessageData struct {
	sender  string
	message string
	color   color.Color
	time    time.Time
}

var chatHistory []chatMessageData
var numLines int
var chatHistoryLock sync.Mutex
var maxShowLines int

const maxHistory = 100

func addToChat(msg irc.ChatMessage) {
	chatHistoryLock.Lock()
	defer chatHistoryLock.Unlock()

	c, err := Hex2Color(msg.Sender.Color)

	var msgColor = color.RGBA{R: 255, G: 255, B: 255}
	if err == nil {
		msgColor = c
	}

	message := sclean.StripControlAndSpecial(msg.Text)
	messageLen := len(message)

	if messageLen > 0 {
		chatHistory = append(chatHistory, chatMessageData{sender: msg.Sender.DisplayName, message: msg.Text, color: msgColor, time: time.Now()})
		numLines++

		trimChatHistory()
	}
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
