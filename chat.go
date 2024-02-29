package main

import (
	"goTwitchGame/sclean"
	"strings"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch/irc"
)

type commandData struct {
	sender  string
	command string
	time    time.Time
}

var userCommands []commandData
var userCommandsLock sync.Mutex
var commandCount = 0

const maxCommandHistory = 10000
const maxCommandLen = 100

func handleChat(msg irc.ChatMessage) {

	message := sclean.StripControlAndSpecial(msg.Text)
	command, isCommand := strings.CutPrefix(message, "!")

	if isCommand {
		handleUserCommand(msg, command)
	}

}

func handleUserCommand(msg irc.ChatMessage, command string) {
	cmdLen := len(command)

	if cmdLen == 0 || cmdLen > maxCommandLen {
		return
	}

	userCommandsLock.Lock()
	userCommands = append(userCommands, commandData{sender: msg.Sender.DisplayName, command: command, time: time.Now()})
	commandCount++
	trimUserCommands()
	userCommandsLock.Unlock()
}

func trimUserCommands() {
	//Remove old lines until we fit
	for {
		if commandCount > maxCommandHistory {
			userCommands = userCommands[1:]
			commandCount--
		} else {
			break
		}
	}
}
