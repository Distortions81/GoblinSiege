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
		Handle: clearVotes,
	},
}

func handleModCommands(msg irc.ChatMessage, command string) bool {

	if msg.Sender.IsBroadcaster || msg.Sender.IsModerator {
		for _, item := range modCommands {
			if strings.EqualFold(item.Name, command) {
				UserMsgDict.Lock.Lock()
				item.Handle()
				UserMsgDict.Lock.Unlock()
				return true
			}
		}
	}

	return false
}

func clearVotes() {
	if UserMsgDict.Count > 0 {
		log.Println("Clearing votes...")
		UserMsgDict.Count = 0
		UserMsgDict.Voting = false
		UserMsgDict.Result = xyi{}
		UserMsgDict.Users = make(map[int64]*userMsgData)
	}
}

func startVote() {
	log.Println("Starting new vote...")
	UserMsgDict.StartTime = time.Now()
	UserMsgDict.Voting = true
	UserMsgDict.Count = 0
	UserMsgDict.Result = xyi{}
	UserMsgDict.Users = make(map[int64]*userMsgData)
}

func endVote() {

	processUserDict()
	if UserMsgDict.Voting {
		log.Println("Ending vote...")
		UserMsgDict.Voting = false
		UserMsgDict.StartTime = time.Now()

		gameMapLock.Lock()
		if UserMsgDict.Count > 0 {
			gameMap[UserMsgDict.Result] = &objectData{Pos: UserMsgDict.Result}
		}
		gameMapLock.Unlock()
	}

}
