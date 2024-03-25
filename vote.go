package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

const maxVoteLen = 10

type VOTE_STATE int

const (
	VOTE_NONE VOTE_STATE = iota
	VOTE_PLAYERS
	VOTE_PLAYERS_DONE
	VOTE_COMPUTER
	VOTE_COMPUTER_DONE
)

type xyi struct {
	X, Y int
}

type xyf32 struct {
	X, Y float32
}

type xyf64 struct {
	X, Y float64
}

type userMsgData struct {
	sender     string
	pos        xyi
	playerMove int
	time       time.Time
}

type userMsgDictData struct {
	Users     map[int64]*userMsgData
	VoteCount int
	VoteState VOTE_STATE

	GameRunning bool
	StartTime   time.Time
	CpuTime     time.Time
	RoundTime   time.Time
	Result      xyi
}

var votes userMsgDictData
var newVoteNotice []*userMsgData

func handleVoteMsg(msg twitch.PrivateMessage, command string) {
	msgLen := len(command)

	//If the message is empty, or huge just discard it
	if msgLen == 0 || msgLen > maxVoteLen {
		return
	}

	args := strings.Split(strings.ToUpper(command), ",")
	numArgs := len(args)
	if numArgs == 2 {

		x, erra := strconv.ParseInt(args[0], 10, 64)
		y, errb := strconv.ParseInt(args[1], 10, 64)

		if erra != nil || errb != nil ||
			x <= 0 || x > boardSizeX ||
			y <= 0 || y > boardSizeY {
			return
		}

		userid := strToID(msg.User.ID)

		userVote := &userMsgData{sender: msg.User.DisplayName, pos: xyi{X: int(x), Y: int(y)}, time: time.Now(), playerMove: board.playerMoveNum}
		votes.Users[userid] = userVote

		if votes.Users[userid] == nil {
			votes.VoteCount++
			newVoteNotice = append(newVoteNotice, userVote)
		}
	}
}

// Currently averages votes
func processVotes() {

	var tX, tY, count uint64

	for _, user := range votes.Users {
		tX += uint64(user.pos.X)
		tY += uint64(user.pos.Y)
		count++
	}
	if count > 0 {
		votes.VoteCount = int(count)
		votes.Result = xyi{X: int(tX / count), Y: int(tY / count)}
	}
	board.playerMoveNum++
}

func setUserScore(id int64, score int) {

	pLock.Lock()
	defer pLock.Unlock()

	players.idmap[id] = &playerData{Points: score}
	players.dirty = true
}

func getUserScore(id int64) (score int, found bool) {

	pLock.Lock()
	defer pLock.Unlock()

	if players.idmap[id] != nil {
		return players.idmap[id].Points, true
	}
	return 0, false
}
