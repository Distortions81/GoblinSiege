package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ServerRunning   bool = true
	ServerIsStopped bool
	roundTime       time.Duration = time.Second * 15
	restTime        time.Duration = time.Second * 5

	skipTwitch *bool
)

func main() {
	skipTwitch = flag.Bool("skip", false, "don't connect to twitch")
	flag.Parse()

	go startEbiten()

	readPlayers()

	players.lock.Lock()
	writePlayers() //Unlocks after serialize

	UserMsgDict.Users = make(map[int64]*userMsgData)
	gameMap = make(map[xyi]*objectData)

	if !*skipTwitch {
		connectTwitch()
	}
	go dbAutoSave()

	startGame()

	//After starting loops, wait here for process signals
	signalHandle := make(chan os.Signal, 1)

	signal.Notify(signalHandle, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-signalHandle

	ServerRunning = false

	qlog("Saving DB...")
	players.lock.Lock()
	writePlayers()
}

func dbAutoSave() {
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
