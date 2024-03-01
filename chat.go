package main

import (
	"goTwitchGame/sclean"
	"strings"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch/irc"
)

const maxCommandLen = 100

type commandData struct {
	sender  string
	command string
	time    time.Time
}

type userMsgDictData struct {
	Users map[int64]*commandData

	Count     int
	Enabled   bool
	StartTime time.Time
	Lock      sync.Mutex
}

var UserMsgDict userMsgDictData

func handleChat(msg irc.ChatMessage) {

	message := sclean.StripControlAndSpecial(msg.Text)
	command, isCommand := strings.CutPrefix(message, "!")

	if isCommand && UserMsgDict.Enabled {
		handleUserCommand(msg, command)
	}
}

func handleUserCommand(msg irc.ChatMessage, command string) {
	cmdLen := len(command)

	if cmdLen == 0 || cmdLen > maxCommandLen {
		return
	}

	UserMsgDict.Lock.Lock()
	if UserMsgDict.Users[msg.Sender.ID] == nil {
		UserMsgDict.Count++
	}
	UserMsgDict.Users[msg.Sender.ID] = &commandData{sender: msg.Sender.DisplayName, command: command, time: time.Now()}
	UserMsgDict.Lock.Unlock()
}

func clearVotes() {
	UserMsgDict.Lock.Lock()
	UserMsgDict.Count = 0
	UserMsgDict.Users = make(map[int64]*commandData)
	UserMsgDict.Lock.Unlock()
}

func startVote() {
	UserMsgDict.Lock.Lock()
	UserMsgDict.StartTime = time.Now()
	UserMsgDict.Enabled = true
	clearVotes()
	UserMsgDict.Lock.Unlock()
}

func endVote() {
	UserMsgDict.Lock.Lock()
	UserMsgDict.Enabled = false
	UserMsgDict.Lock.Unlock()
}
