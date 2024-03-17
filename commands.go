package main

import (
	"fmt"
	"math/rand"
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
		if cmd.Name == "help" || cmd.Name == "modHelp" {
			continue
		}
		modCmdHelp = append(modCmdHelp, commandData{Name: cmd.Name, Desc: cmd.Desc})
	}
}

var modCommands []commandData = []commandData{
	{
		Name:   "help",
		Desc:   "Show player commands.",
		Handle: helpCommand,
	},
	{
		Name:   "modHelp",
		Desc:   "Show moderator commands.",
		Handle: modHelpCommand,
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
		Desc:   "Manually start a new voting round.",
		Handle: startVote,
	},
	{
		Name:   "endVote",
		Desc:   "Manually end a voting round.",
		Handle: endVote,
	},
	{
		Name:   "clearVotes",
		Desc:   "clear all votes.",
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

func clearGameBoard() {
	qlog("Clearing game board...")
	board.lock.Lock()

	board.bmap = make(map[xyi]*objectData)
	tower1 := getOtype("Stone Tower")
	for i := 0; i < 10; i++ {
		tpos := xyi{X: rand.Intn(boardSizeX-1) + 1, Y: rand.Intn(boardSizeY-1) + 1}
		board.bmap[tpos] = &objectData{Pos: tpos, oTypeP: tower1, Health: 100}
	}
	board.lock.Unlock()
}

func clearVotes() {
	qlog("Resetting votes...")
	UserMsgDict.VoteCount = 0
	UserMsgDict.Result = xyi{}
	UserMsgDict.Users = make(map[int64]*userMsgData)
}

func startGame() {
	qlog("Starting game...")

	clearGameBoard()

	UserMsgDict.GameRunning = true

	startVote()
}

func endGame() {
	qlog("Stopping game...")

	clearGameBoard()
	clearVotes()

	UserMsgDict.VoteState = VOTE_NONE
	UserMsgDict.GameRunning = false

}

func startVote() {
	clearVotes()
	qlog("Starting new vote...")
	UserMsgDict.StartTime = time.Now()
	UserMsgDict.VoteState = VOTE_PLAYERS
}

// Locks and unlocks board
func endVote() {

	processUserDict()
	if UserMsgDict.VoteState == VOTE_PLAYERS {
		qlog("Ending vote...")
		UserMsgDict.VoteState = VOTE_PLAYERS_DONE
		UserMsgDict.StartTime = time.Now()

		if UserMsgDict.VoteCount > 0 &&
			board.bmap[UserMsgDict.Result] == nil &&
			UserMsgDict.Result.X > 0 &&
			UserMsgDict.Result.X <= boardSizeX &&
			UserMsgDict.Result.Y > 0 &&
			UserMsgDict.Result.Y <= boardSizeY {

			board.lock.Lock()
			board.bmap[UserMsgDict.Result] = &objectData{Pos: UserMsgDict.Result}
			board.lock.Unlock()
		}
	}
}

func modHelpCommand() {
	var buf string
	for _, cmd := range modCmdHelp {
		buf = buf + fmt.Sprintf("%v%v -- %v, ", userSettings.CmdPrefix, cmd.Name, cmd.Desc)
	}
	fSay(buf)
}

func helpCommand() {
	buf := fmt.Sprintf("%vx,y", userSettings.CmdPrefix)
	fSay(buf)
}
