package main

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

const maxCommandLen = 100

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
	Count       int
	Voting      bool
	GameStarted bool
	StartTime   time.Time
	Lock        sync.Mutex
	RoundTime   time.Time
	Result      xyi
}

var (
	UserMsgDict userMsgDictData
)

func handleUserDictMsg(msg twitch.PrivateMessage, command string) {
	msgLen := len(command)

	if msgLen == 0 || msgLen > maxCommandLen {
		return
	}

	userid := strToID(msg.User.ID)

	UserMsgDict.Lock.Lock()
	if UserMsgDict.Users[userid] == nil {
		UserMsgDict.Count++
	}
	UserMsgDict.Users[userid] = &userMsgData{sender: msg.User.DisplayName, command: command, time: time.Now()}
	UserMsgDict.Lock.Unlock()
}

func processUserDict() {

	var tX, tY, count uint64

	for _, user := range UserMsgDict.Users {
		args := strings.Split(strings.ToUpper(user.command), " ")
		if len(args) > 1 {

			//Convert from text to value
			xb := []byte(args[0])
			if len(xb) == 0 {
				break
			}
			x := int(xb[0]) - 64
			y, err := strconv.ParseInt(args[1], 10, 64)

			qlog("user: %v, x: %v, y: %v", user.sender, x, y)
			if err != nil ||
				x < 1 || x > boardSize ||
				y < 1 || y > boardSize {
				continue
			}

			tX += uint64(x)
			tY += uint64(y)

			count++
		}

	}
	if count > 0 {
		UserMsgDict.Count = int(count)
		UserMsgDict.Result = xyi{X: int(tX / count), Y: int(tY / count)}
	} else {
		qlog("Not enough votes.")
	}
}

func setUserScore(id int64, score int) {

	dbLock.Lock()
	defer dbLock.Unlock()

	Players[id] = &playerData{Points: score}
	dbDirty = true
}

func getUserScore(id int64) (score int, found bool) {

	dbLock.Lock()
	defer dbLock.Unlock()

	if Players[id] != nil {
		return Players[id].Points, true
	}
	return 0, false
}
