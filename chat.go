package main

import (
	"goTwitchGame/sclean"
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

func handleChat(msg irc.ChatMessage) {

	message := sclean.StripControlAndSpecial(msg.Text)
	command, isCommand := strings.CutPrefix(message, "!")

	if isCommand && UserMsgDict.Enabled {
		handleUserDictMsg(msg, command)
	}
}

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

func clearVotes() {
	UserMsgDict.Lock.Lock()
	UserMsgDict.Count = 0
	UserMsgDict.Users = make(map[int64]*userMsgData)
	UserMsgDict.Lock.Unlock()
}

func startVote() {
	UserMsgDict.Lock.Lock()
	UserMsgDict.StartTime = time.Now()
	UserMsgDict.Enabled = true
	UserMsgDict.Count = 0
	UserMsgDict.Users = make(map[int64]*userMsgData)
	UserMsgDict.Lock.Unlock()
}

func endVote() {
	log.Println("Ending vote...")
	UserMsgDict.Lock.Lock()
	UserMsgDict.Enabled = false
	UserMsgDict.Lock.Unlock()

	processUserDict()
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
