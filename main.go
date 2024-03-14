package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ServerRunning bool          = true
	roundTime     time.Duration = time.Second * 15
	restTime      time.Duration = time.Second * 5

	skipTwitch *bool
)

func main() {
	skipTwitch = flag.Bool("skip", false, "don't connect to twitch")
	flag.Parse()

	go startEbiten()

	readPlayers()

	UserMsgDict.Users = make(map[int64]*userMsgData)

	readSettings()
	if !*skipTwitch {
		connectTwitch()
	}
	go playersAutosave()

	board.bmap = make(map[xyi]*objectData)
	board.bmap[UserMsgDict.Result] = &objectData{Pos: xyi{X: 1, Y: 1}}
	startGame()

	//After starting loops, wait here for process signals
	signalHandle := make(chan os.Signal, 1)

	signal.Notify(signalHandle, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-signalHandle

	ServerRunning = false

	qlog("Saving players...")
	players.lock.Lock()
	writePlayers()
}

func playersAutosave() {
	for ServerRunning {

		players.lock.Lock()
		if players.dirty {
			players.dirty = false
			writePlayers() //This unlocks after serialize
		} else {
			//No write to do, unlock
			players.lock.Unlock()
		}

		time.Sleep(time.Second * 30)
	}
}
