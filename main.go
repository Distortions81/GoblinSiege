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
	playerRoundTime time.Duration = time.Second * 10
	cpuRoundTime    time.Duration = time.Second

	skipTwitch *bool
	debugMode  *bool
)

func main() {
	skipTwitch = flag.Bool("skip", false, "don't connect to twitch")
	debugMode = flag.Bool("debug", false, "debug mode")
	flag.Parse()

	if *debugMode {
		cpuRoundTime = time.Microsecond
		playerRoundTime = time.Millisecond * 100
	}

	//Wait here for process signals
	signalHandle := make(chan os.Signal, 1)

	//Start ebiten game lib
	go func() {
		startEbiten()

		//Exit if window closed
		signalHandle <- os.Interrupt
	}()

	//Read player scores
	readPlayers()

	//Load settings
	readSettings()

	//Connect to twitch
	if !*skipTwitch {
		connectTwitch()
	}

	//Start autosave loop
	go playersAutosave()

	//Voting loop
	go handleRounds()

	//Start the game mode
	startGame()

	signal.Notify(signalHandle, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-signalHandle

	//Shutdown server and save
	ServerRunning = false

	players.lock.Lock()
	writePlayers()
}
