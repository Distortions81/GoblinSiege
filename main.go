package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	MODE_SPLASH = iota
	MODE_PLAY_TWITCH
	MODE_PLAY_SINGLE
	MODE_SETTINGS
)

var (
	ServerRunning  atomic.Bool
	gameMode       atomic.Int32
	playerMoveTime time.Duration = time.Second * 10
	cpuMoveTime    time.Duration = time.Second * 2

	skipTwitch, fastMode, noTowers, smartMove, debugMode, skipMenu *bool

	gameLoaded   atomic.Bool
	signalHandle chan os.Signal
)

func main() {
	skipTwitch = flag.Bool("skipTwitch", false, "don't connect to twitch")
	fastMode = flag.Bool("fast", false, "fast mode")
	noTowers = flag.Bool("notower", false, "don't spawn towers")
	smartMove = flag.Bool("smartmove", false, "Use intelligent moves to simulate a coordinated audiance.")
	debugMode = flag.Bool("debug", false, "print debug info")
	skipMenu = flag.Bool("skipMenu", false, "Skips main menu")
	flag.Parse()

	ServerRunning.Store(true)

	// After starting loops, wait here for process signals
	signalHandle = make(chan os.Signal, 1)
	signal.Notify(signalHandle, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	var err error
	splash, _, err = ebitenutil.NewImageFromFile("data/sprites/splash.png")
	if err != nil {
		log.Fatal(err)
	}

	go startEbiten()

	loadAssets()

	readSettings()

	if !*skipTwitch {
		readPlayers()
		connectTwitch()
		go playersAutosave()
	}

	go handleMoves()

	go aniTimer()

	board.towerMap = make(map[xyi]*objectData)
	board.goblinMap = make(map[xyi]*objectData)
	votes.Users = make(map[int64]*userMsgData)

	board.fFrame = ebiten.NewImage(defaultWindowWidth, defaultWindowHeight)
	board.checkerCache = ebiten.NewImage(defaultWindowWidth, defaultWindowHeight)
	board.deadCache = ebiten.NewImage(defaultWindowWidth, defaultWindowHeight)
	board.checkerDirty = true

	setupButtons()

	gameLoaded.Store(true)

	if *skipMenu {
		gameMode.Store(MODE_PLAY_TWITCH)
	}

	<-signalHandle

	ServerRunning.Store(false)
	writePlayers()
	time.Sleep(time.Second)
}

// Used for action animations
func aniTimer() {
	for ServerRunning.Load() {
		aniCount.Add(1)
		time.Sleep(time.Second / 3)
	}
}
