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

var userCommands map[int64]*commandData
var userCommandsLock sync.Mutex
var commandCount = 0
var votesEnabled bool = true
var voteStarted time.Time

const maxCommandLen = 100

func handleChat(msg irc.ChatMessage) {

	message := sclean.StripControlAndSpecial(msg.Text)
	command, isCommand := strings.CutPrefix(message, "!")

	if isCommand && votesEnabled {
		handleUserCommand(msg, command)
	}
}

func handleUserCommand(msg irc.ChatMessage, command string) {
	cmdLen := len(command)

	if cmdLen == 0 || cmdLen > maxCommandLen {
		return
	}

	userCommandsLock.Lock()
	if userCommands[msg.Sender.ID] == nil {
		commandCount++
	}
	userCommands[msg.Sender.ID] = &commandData{sender: msg.Sender.DisplayName, command: command, time: time.Now()}
	userCommandsLock.Unlock()
}

func clearVotes() {
	userCommandsLock.Lock()
	commandCount = 0
	userCommands = make(map[int64]*commandData)
	userCommandsLock.Unlock()
}

func startVote() {
	userCommandsLock.Lock()
	voteStarted = time.Now()
	votesEnabled = true
	clearVotes()
	userCommandsLock.Unlock()
}

func endVote() {
	userCommandsLock.Lock()
	votesEnabled = false
	userCommandsLock.Unlock()
}
