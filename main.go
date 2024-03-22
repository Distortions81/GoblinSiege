package main

import (
	"flag"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ServerRunning  bool          = true
	playerMoveTime time.Duration = time.Second * 10
	cpuMoveTime    time.Duration = time.Second * 2
	maxMoves                     = 100

	skipTwitch, fastMode, noTowers, smartMove, debugMode *bool
)

func main() {
	skipTwitch = flag.Bool("skip", false, "don't connect to twitch")
	fastMode = flag.Bool("fast", false, "fast mode")
	noTowers = flag.Bool("notower", false, "don't spawn towers")
	smartMove = flag.Bool("smartmove", false, "Use intelligent moves to simulate a coordinated audiance.")
	debugMode = flag.Bool("debug", false, "print debug info")
	flag.Parse()

	board.towerMap = make(map[xyi]*objectData)
	board.goblinMap = make(map[xyi]*objectData)
	votes.Users = make(map[int64]*userMsgData)

	freezeFrame = ebiten.NewImage(defaultWindowWidth, defaultWindowHeight)

	loadAssets()

	go aniTimer()

	readPlayers()

	readSettings()

	if !*skipTwitch {
		connectTwitch()
	}

	go playersAutosave()

	go handleMoves()

	startGame()

	startEbiten() //Blocks until exit

	ServerRunning = false
	writePlayers()
}

// Used for action animations
func aniTimer() {
	for {
		aniCount++
		time.Sleep(time.Second / 3)
	}
}
