package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ServerRunning   bool          = true
	playerRoundTime time.Duration = time.Second * 15
	cpuRoundTime    time.Duration = time.Second * 5

	skipTwitch *bool
)

func main() {
	skipTwitch = flag.Bool("skip", false, "don't connect to twitch")
	flag.Parse()

	//Start ebiten game lib
	go startEbiten()

	//Read player scores
	readPlayers()

	//Load settings
	readSettings()

	//Connect to twitch
	if !*skipTwitch {
		connectTwitch()
	}

	//Start autosave loop (replace later)
	go playersAutosave()

	//Start the game mode
	startGame()

	//Wait here for process signals
	signalHandle := make(chan os.Signal, 1)

	signal.Notify(signalHandle, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-signalHandle

	//Shutdown server and save
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
