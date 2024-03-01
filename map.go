package main

import (
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type objectData struct {
	Pos  xyi
	Type int
}

const mag = 24
const size = 23
const boardSize = 25
const itemOffset = 2
const boardPixels = boardSize * mag

var (
	gameMap      map[xyi]*objectData
	gameMapLock  sync.Mutex
	gameMapCount int
)

func drawGameBoard(screen *ebiten.Image) {

	gameMapLock.Lock()
	defer gameMapLock.Unlock()

	if gameMapCount == 0 {
		return
	}

	vector.DrawFilledRect(screen, size*itemOffset, size*itemOffset, boardPixels, boardPixels, ColorDarkGreen, true)
	for _, item := range gameMap {
		vector.DrawFilledCircle(screen, float32((item.Pos.X+itemOffset)*mag), float32((item.Pos.Y+itemOffset)*mag), size, color.White, true)
	}
}
