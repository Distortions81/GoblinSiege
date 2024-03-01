package main

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch/irc"
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

	Count     int
	Enabled   bool
	StartTime time.Time
	Lock      sync.Mutex
	RoundTime time.Time
	Result    xyi
}

var (
	UserMsgDict userMsgDictData
)

func handleUserDictMsg(msg irc.ChatMessage, command string) {
	msgLen := len(command)

	if msgLen == 0 || msgLen > maxCommandLen {
		return
	}

	UserMsgDict.Lock.Lock()
	if UserMsgDict.Users[msg.Sender.ID] == nil {
		UserMsgDict.Count++
	}
	UserMsgDict.Users[msg.Sender.ID] = &userMsgData{sender: msg.Sender.DisplayName, command: command, time: time.Now()}
	UserMsgDict.Lock.Unlock()
}

func processUserDict() {

	var tX, tY, count uint64

	for _, user := range UserMsgDict.Users {
		args := strings.Split(strings.ToUpper(user.command), " ")
		if len(args) > 1 {

			//Convert from text to value
			xb := []byte(args[0])
			x := int(xb[0]) - 65
			y, err := strconv.ParseInt(args[1], 10, 64)

			if err != nil || x < 0 || x > 90 {
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
		log.Printf("Average from %v users: %v,%v\n", count, UserMsgDict.Result.X, UserMsgDict.Result.Y)
	} else {
		log.Println("Not enough votes.")
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
