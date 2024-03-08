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

const size = 23
const mag = size + 1
const boardSize = 25
const itemOffset = 2
const boardPixels = ((boardSize + 1) * mag) + 2

var (
	gameMap     map[xyi]*objectData
	gameMapLock sync.Mutex
)

func drawGameBoard(screen *ebiten.Image) {

	gameMapLock.Lock()
	defer gameMapLock.Unlock()

	vector.DrawFilledRect(screen, size*itemOffset, size*itemOffset, boardPixels, boardPixels, ColorVeryDarkGreen, true)
	for _, item := range gameMap {
		vector.DrawFilledCircle(screen, float32((item.Pos.X+itemOffset)*mag), float32((item.Pos.Y+itemOffset)*mag), size, color.White, true)
	}
}
