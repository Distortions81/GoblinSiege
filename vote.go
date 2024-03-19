package main

import (
	"strconv"
	"strings"
	"sync"
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

type userMsgData struct {
	sender string
	pos    xyi
	time   time.Time
}

type userMsgDictData struct {
	Users map[int64]*userMsgData

	GameRunning bool

	VoteCount int
	VoteState VOTE_STATE
	StartTime time.Time
	CpuTime   time.Time
	Lock      sync.Mutex
	RoundTime time.Time
	Result    xyi
}

var (
	votes userMsgDictData
)

func handleVoteMsg(msg twitch.PrivateMessage, command string) {
	msgLen := len(command)

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
		votes.Lock.Lock()
		if votes.Users[userid] == nil {
			votes.VoteCount++
		}
		votes.Users[userid] = &userMsgData{sender: msg.User.DisplayName, pos: xyi{X: int(x), Y: int(y)}, time: time.Now()}

		votes.Lock.Unlock()

	}

}

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
}

func setUserScore(id int64, score int) {

	players.lock.Lock()
	defer players.lock.Unlock()

	players.idmap[id] = &playerData{Points: score}
	players.dirty = true
}

func getUserScore(id int64) (score int, found bool) {

	players.lock.Lock()
	defer players.lock.Unlock()

	if players.idmap[id] != nil {
		return players.idmap[id].Points, true
	}
	return 0, false
}
