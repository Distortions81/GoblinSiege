package main

import "time"

const (
	defaultWindowWidth  = 1280
	defaultWindowHeight = 720
	deathDelay          = time.Millisecond * 300

	playersFile = "players.json"
	authFile    = "settings.json"

	size         = 32
	mag          = size
	boardSizeX   = 20
	boardSizeY   = 20
	enemyBoardX  = 15
	offX         = 5
	offY         = 1
	offPixX      = size * offX
	offPixY      = size * offY
	boardPixelsX = ((boardSizeX) * mag)
	boardPixelsY = ((boardSizeY) * mag)

	defaultVolume = 0.5
)
