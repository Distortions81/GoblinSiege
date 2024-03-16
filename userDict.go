package main

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

const maxDictMsgLen = 10

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
	sender  string
	command string
	time    time.Time
}

type userMsgDictData struct {
	Users map[int64]*userMsgData

	GameRunning bool

	VoteCount int
	VoteState VOTE_STATE
	StartTime time.Time
	Lock      sync.Mutex
	RoundTime time.Time
	Result    xyi
}

var (
	UserMsgDict userMsgDictData
)

func handleUserDictMsg(msg twitch.PrivateMessage, command string) {
	msgLen := len(command)

	if msgLen == 0 || msgLen > maxDictMsgLen {
		return
	}

	//Move this later, we only need to do on score
	userid := strToID(msg.User.ID)
	UserMsgDict.Lock.Lock()
	if UserMsgDict.Users[userid] == nil {
		UserMsgDict.VoteCount++
	}
	UserMsgDict.Users[userid] = &userMsgData{sender: msg.User.DisplayName, command: command, time: time.Now()}
	UserMsgDict.Lock.Unlock()
}

func processUserDict() {

	var tX, tY, count uint64

	for _, user := range UserMsgDict.Users {
		args := strings.Split(strings.ToUpper(user.command), ",")
		numArgs := len(args)
		if numArgs == 2 {

			x, erra := strconv.ParseInt(args[0], 10, 64)
			y, errb := strconv.ParseInt(args[1], 10, 64)

			qlog("user: %v, x: %v, y: %v", user.sender, x, y)
			if erra != nil || errb != nil ||
				x < 1 || x > boardSizeX ||
				y < 1 || y > boardSizeY {
				continue
			}

			tX += uint64(x)
			tY += uint64(y)

			count++
		}

	}
	if count > 0 {
		UserMsgDict.VoteCount = int(count)
		UserMsgDict.Result = xyi{X: int(tX / count), Y: int(tY / count)}
	} else {
		qlog("Not enough votes.")
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
