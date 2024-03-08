package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
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
		if cmd.Name == "modHelp" {
			continue
		}
		modCmdHelp = append(modCmdHelp, commandData{Name: cmd.Name, Desc: cmd.Desc})
	}
}

var modCommands []commandData = []commandData{
	{
		Name:   "modHelp",
		Handle: helpCommand,
	},
	{
		Name:   "startGame",
		Desc:   "Start game, voting will start and end automatically",
		Handle: startGame,
	},
	{
		Name:   "endGame",
		Desc:   "End game and clear votes.",
		Handle: endGame,
	},
	{
		Name:   "startVote",
		Desc:   "Manually start a voting round",
		Handle: startVote,
	},
	{
		Name:   "endVote",
		Desc:   "Manually end a voting round",
		Handle: endVote,
	},
	{
		Name:   "clearVotes",
		Desc:   "Stop and clear all votes",
		Handle: clearVotes,
	},
}

func handleModCommands(msg twitch.PrivateMessage, command string) bool {

	if strings.EqualFold(msg.User.Name, userSettings.UserName) {
		for _, item := range modCommands {
			if strings.EqualFold(item.Name, command) {
				if item.Handle == nil {
					sayLog("Command %v%v has nil func.", userSettings.CmdPrefix, item.Name)
					continue
				}
				UserMsgDict.Lock.Lock()
				item.Handle()
				UserMsgDict.Lock.Unlock()
				return true
			}
		}
	}

	return false
}

func startGame() {
	qlog("Starting game...")

	//Clear game board
	gameMapLock.Lock()
	gameMap = make(map[xyi]*objectData)
	gameMapLock.Unlock()

	//Reset votes
	UserMsgDict.Count = 0
	UserMsgDict.Voting = false
	UserMsgDict.Result = xyi{}
	UserMsgDict.Users = make(map[int64]*userMsgData)

	UserMsgDict.GameRunning = true

	startVote()
}

func endGame() {
	qlog("Stopping game...")

	UserMsgDict.Count = 0
	UserMsgDict.Voting = false
	UserMsgDict.Result = xyi{}
	UserMsgDict.Users = make(map[int64]*userMsgData)
	UserMsgDict.GameRunning = false

}

func clearVotes() {
	if UserMsgDict.Count > 0 {
		qlog("Clearing votes...")
		UserMsgDict.Count = 0
		UserMsgDict.Voting = false
		UserMsgDict.Result = xyi{}
		UserMsgDict.Users = make(map[int64]*userMsgData)
	}
}

func startVote() {
	qlog("Starting new vote...")
	UserMsgDict.StartTime = time.Now()
	UserMsgDict.Voting = true
	UserMsgDict.Count = 0
	UserMsgDict.Result = xyi{}
	UserMsgDict.Users = make(map[int64]*userMsgData)
}

func endVote() {

	processUserDict()
	if UserMsgDict.Voting {
		qlog("Ending vote...")
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
	var buf string
	for _, cmd := range modCmdHelp {
		buf = buf + fmt.Sprintf("%v%v -- %v, ", userSettings.CmdPrefix, cmd.Name, cmd.Desc)
	}
	fSay(buf)
}
