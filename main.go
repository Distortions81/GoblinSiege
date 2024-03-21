package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ServerRunning  bool          = true
	playerMoveTime time.Duration = time.Second * 10
	cpuMoveTime    time.Duration = time.Second * 2
	maxMoves                     = 100

	skipTwitch *bool
	fastMode   *bool
)

func main() {
	skipTwitch = flag.Bool("skip", false, "don't connect to twitch")
	fastMode = flag.Bool("fast", false, "fast mode")
	flag.Parse()

	board.playMap = make(map[xyi]*objectData)
	board.enemyMap = make(map[xyi]*objectData)
	votes.Users = make(map[int64]*userMsgData)

	//Wait here for process signals
	signalHandle := make(chan os.Signal, 1)

	loadEmbed()

	go aniTimer()

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
	go handleMoves()

	//Start the game mode
	startGame()

	//Start ebiten
	startEbiten()

	signal.Notify(signalHandle, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-signalHandle

	//Shutdown server and save
	ServerRunning = false

	players.lock.Lock()
	writePlayers()
}

func aniTimer() {
	for {
		aniCount++
		time.Sleep(time.Second / 3)
	}
}
