package main

import (
	"sync"
)

var (
	updateGameSizeLock sync.Mutex

	ScreenWidth, ScreenHeight int
)
