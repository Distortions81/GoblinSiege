package main

import (
	"flag"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	MODE_SPLASH = iota
	MODE_PLAY_TWITCH
	MODE_PLAY_SINGLE
	MODE_SETTINGS
)

var (
	ServerRunning  bool = true
	gameMode       int
	playerMoveTime time.Duration = time.Second * 10
	cpuMoveTime    time.Duration = time.Second * 2

	skipTwitch, fastMode, noTowers, smartMove, debugMode, skipMenu *bool
)

func main() {
	skipTwitch = flag.Bool("skipTwitch", false, "don't connect to twitch")
	fastMode = flag.Bool("fast", false, "fast mode")
	noTowers = flag.Bool("notower", false, "don't spawn towers")
	smartMove = flag.Bool("smartmove", false, "Use intelligent moves to simulate a coordinated audiance.")
	debugMode = flag.Bool("debug", false, "print debug info")
	skipMenu = flag.Bool("skipMenu", false, "Skips main menu")
	flag.Parse()

	if *skipMenu {
		gameMode = MODE_PLAY_TWITCH
	}

	board.towerMap = make(map[xyi]*objectData)
	board.goblinMap = make(map[xyi]*objectData)
	votes.Users = make(map[int64]*userMsgData)

	board.fFrame = ebiten.NewImage(defaultWindowWidth, defaultWindowHeight)
	board.checkerCache = ebiten.NewImage(defaultWindowWidth, defaultWindowHeight)
	board.deadCache = ebiten.NewImage(defaultWindowWidth, defaultWindowHeight)
	board.checkerDirty = true

	loadAssets()

	go aniTimer()

	readPlayers()

	readSettings()

	if !*skipTwitch {
		connectTwitch()
	}

	go playersAutosave()

	go handleMoves()

	startEbiten() //Blocks until exit

	ServerRunning = false
	writePlayers()
}

// Used for action animations
func aniTimer() {
	for {
		aniCount.Add(1)
		time.Sleep(time.Second / 3)
	}
}
