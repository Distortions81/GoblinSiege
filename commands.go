package main

import (
	"log"
	"strings"
	"time"

	"github.com/Adeithe/go-twitch/irc"
)

type commandData struct {
	Name   string
	Desc   string
	Handle func()
}

/* Lame workaround for initialization cycle error */
var modCmdHelp []commandData

func init() {
	for _, cmd := range modCommands {
		modCmdHelp = append(modCmdHelp, commandData{Name: cmd.Name, Desc: cmd.Desc})
	}
}

var modCommands []commandData = []commandData{
	{
		Name:   "help",
		Desc:   "(You are here)",
		Handle: helpCommand,
	},
	{
		Name:   "startGame",
		Desc:   "Start game, voting will start and end automatically.",
		Handle: nil,
	},
	{
		Name:   "startVote",
		Desc:   "Manually start a voting round.",
		Handle: startVote,
	},
	{
		Name:   "endVote",
		Desc:   "Manually end a voting round.",
		Handle: endVote,
	},
	{
		Name:   "clearVotes",
		Desc:   "Stop and clear all votes.",
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

func helpCommand() {
	for _, cmd := range modCmdHelp {
		log.Printf("!%v -- %v\n", cmd.Name, cmd.Desc)
	}
}
