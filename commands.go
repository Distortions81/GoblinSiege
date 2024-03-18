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
				votes.Lock.Lock()
				item.Handle()
				votes.Lock.Unlock()
				return true
			}
		}
	}

	return false
}

func clearGameBoard() {
	qlog("Clearing game board...")
	board.lock.Lock()

	board.playMap = make(map[xyi]*objectData)
	board.enemyMap = make(map[xyi]*objectData)
	board.roundNum = 0
	board.arrowsShot = make([]arrowData, 0)

	board.lock.Unlock()
}

func clearVotes() {
	qlog("Resetting votes...")
	votes.VoteCount = 0
	votes.Result = xyi{}
	votes.Users = make(map[int64]*userMsgData)
}

func startGame() {
	qlog("Starting game...")

	clearGameBoard()

	board.gameover = GAME_RUNNING
	votes.VoteCount = 0
	votes.GameRunning = true

	startVote()
}

func endGame() {
	qlog("Stopping game...")
	clearVotes()

	votes.VoteState = VOTE_NONE
	votes.GameRunning = false

}

func startVote() {
	qlog("Starting new vote...")
	votes.StartTime = time.Now()
	votes.VoteState = VOTE_PLAYERS
}

// Locks and unlocks board
func endVote() {

	processUserDict()
	board.lock.Lock()
	defer board.lock.Unlock()

	if votes.VoteState == VOTE_PLAYERS {
		qlog("Ending vote...")

		addTower()
		votes.VoteState = VOTE_PLAYERS_DONE
		votes.StartTime = time.Now()
	}
	clearVotes()
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

func init() {
	for _, cmd := range modCommands {
		if cmd.Name == "help" || cmd.Name == "modHelp" {
			continue
		}
		modCmdHelp = append(modCmdHelp, commandData{Name: cmd.Name, Desc: cmd.Desc})
	}
}
