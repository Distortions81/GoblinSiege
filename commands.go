package main

import (
	"log"
	"strings"
	"time"

	"github.com/Adeithe/go-twitch/irc"
)

type commandData struct {
	Name   string
	Handle func()
}

var modCommands []commandData = []commandData{
	{
		Name:   "start",
		Handle: startVote,
	},
	{
		Name:   "end",
		Handle: endVote,
	},
	{
		Name:   "clear",
		Handle: endVote,
	},
}

func handleModCommands(msg irc.ChatMessage, command string) bool {

	if msg.Sender.IsBroadcaster || msg.Sender.IsModerator {
		for _, item := range modCommands {
			if strings.EqualFold(item.Name, command) {
				item.Handle()
				return true
			}
		}
	}

	return false
}

func clearVotes() {
	UserMsgDict.Lock.Lock()
	if UserMsgDict.Count > 0 {
		log.Println("Clearing votes...")
		UserMsgDict.Count = 0
		UserMsgDict.Result = xyi{}
		UserMsgDict.Users = make(map[int64]*userMsgData)
	}
	UserMsgDict.Lock.Unlock()
}

func startVote() {
	UserMsgDict.Lock.Lock()
	log.Println("Starting new vote...")
	UserMsgDict.StartTime = time.Now()
	UserMsgDict.Enabled = true
	UserMsgDict.Count = 0
	UserMsgDict.Result = xyi{}
	UserMsgDict.Users = make(map[int64]*userMsgData)
	UserMsgDict.Lock.Unlock()
}

func endVote() {

	UserMsgDict.Lock.Lock()
	if UserMsgDict.Enabled {
		log.Println("Ending vote...")
		UserMsgDict.Enabled = false
	}
	UserMsgDict.Lock.Unlock()

	processUserDict()
}
