package main

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	size         = 22
	mag          = size + 1
	boardSize    = 25
	itemOffset   = 2
	offsetPixels = size * itemOffset
	boardPixels  = ((boardSize) * mag)
)

type objectData struct {
	Pos  xyi
	Type int
}

var (
	gameMap     map[xyi]*objectData
	gameMapLock sync.Mutex
)

func drawGameBoard(screen *ebiten.Image) {

	gameMapLock.Lock()
	defer gameMapLock.Unlock()

	//Draw left side bg red
	for x := 0; x < boardSize; x++ {
		for y := 0; y < boardSize; y++ {
			color := ColorDarkRed
			if (x+y)%2 == 0 {
				color = ColorRed
			}
			vector.DrawFilledRect(screen, float32(mag*x)+offsetPixels, float32(mag*y)+offsetPixels, size, size, color, true)
		}
	}

	//Draw right side bg blue
	for x := 0; x < boardSize; x++ {
		for y := 0; y < boardSize; y++ {
			color := ColorBlue
			if (x+y)%2 == 0 {
				color = ColorDarkBlue
			}
			vector.DrawFilledRect(screen, float32(mag*x)+(offsetPixels)+boardPixels, float32(mag*y)+offsetPixels, size, size, color, true)
		}
	}

	//Draw towers
	for _, item := range gameMap {
		vector.DrawFilledCircle(screen, float32((item.Pos.X+itemOffset)*mag)-(size/1.5), float32((item.Pos.Y+itemOffset)*mag)-(size/1.5), size/2, color.White, true)
	}

	buf := fmt.Sprintf("%v", boardSize)
	text.Draw(screen, "X:0", monoFont, 0, 60, color.White)
	text.Draw(screen, buf, monoFont, 0, boardPixels+40, color.White)

	text.Draw(screen, "Y:0", monoFont, 40, 40, color.White)
	text.Draw(screen, buf, monoFont, boardPixels+20, 40, color.White)

	text.Draw(screen, "Audience", monoFont, boardPixels/2, boardPixels+65, color.White)
	text.Draw(screen, "Computer", monoFont, boardPixels+boardPixels/2, boardPixels+65, color.White)
}
