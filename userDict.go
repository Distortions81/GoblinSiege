package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Adeithe/go-twitch/irc"
)

const maxCommandLen = 100

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
}

var UserMsgDict userMsgDictData

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
	UserMsgDict.Lock.Lock()
	defer UserMsgDict.Lock.Unlock()

	var tX, tY, count uint64

	for _, user := range UserMsgDict.Users {
		args := strings.Split(strings.ToUpper(user.command), " ")
		if len(args) > 1 {

			//Convert from text to value
			xb := []byte(args[0])
			yb := []byte(args[1])
			x := int(xb[0]) - 65
			y := int(yb[0]) - 65

			if x < 0 || x > 90 || y < 0 || y > 90 {
				continue
			}

			tX += uint64(x)
			tY += uint64(y)
			count++
		}
	}
	if count > 0 {
		avrX := tX / count
		avrY := tY / count
		log.Printf("Average from %v users: %v,%v\n", count, avrX, avrY)
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
