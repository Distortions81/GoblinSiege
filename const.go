package main

import "time"

const (
	version             = "0.2.29"
	defaultWindowWidth  = 1280
	defaultWindowHeight = 720
	deathDelay          = time.Millisecond * 500
	attackDelay         = time.Millisecond * 100

	playersFile = "players.json"
	authFile    = "settings.json"

	size         = 32
	mag          = size
	hMag         = mag / 2
	boardSizeX   = 20
	boardSizeY   = 20
	enemyBoardX  = 15
	offX         = 5
	offY         = 1
	offPixX      = size * offX
	offPixY      = size * offY
	boardPixelsX = ((boardSizeX) * mag)
	boardPixelsY = ((boardSizeY) * mag)

	defaultVolume   = 0.5
	gameOverFadeSec = 3
	arrowSpeed      = 840
	flashSpeed      = time.Millisecond * 125
)
