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

	skipTwitch, fastMode, noTowers, smartMove *bool
)

func main() {
	skipTwitch = flag.Bool("skip", false, "don't connect to twitch")
	fastMode = flag.Bool("fast", false, "fast mode")
	noTowers = flag.Bool("notower", false, "don't spawn towers")
	smartMove = flag.Bool("smartmove", false, "Use intelligent moves to simulate a coordinated audiance.")
	flag.Parse()

	board.playMap = make(map[xyi]*objectData)
	board.enemyMap = make(map[xyi]*objectData)
	votes.Users = make(map[int64]*userMsgData)

	freezeFrame = ebiten.NewImage(defaultWindowWidth, defaultWindowHeight)

	//Load fonts, sprites and sounds
	loadEmbed()

	//Used for action animations
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

	//Shutdown server and save
	ServerRunning = false

	players.lock.Lock()
	writePlayers()
}

// Used for action animations
// TODO move this to the ebiten update handler
func aniTimer() {
	for {
		aniCount++
		time.Sleep(time.Second / 3)
	}
}
