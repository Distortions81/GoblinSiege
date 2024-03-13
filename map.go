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
	boardCache  *ebiten.Image
	boardCached bool
)

func drawGameBoard(screen *ebiten.Image) {
	updateGameSizeLock.Lock()
	defer updateGameSizeLock.Unlock()

	gameMapLock.Lock()
	defer gameMapLock.Unlock()

	if !boardCached {
		boardCache.Clear()

		//Draw left side bg red
		for x := 0; x < boardSize; x++ {
			for y := 0; y < boardSize; y++ {
				tColor := ColorDarkRed
				if (x+y)%2 == 0 {
					tColor = ColorRed
				}
				vector.DrawFilledRect(boardCache, float32(mag*x)+offsetPixels, float32(mag*y)+offsetPixels, size, size, tColor, true)

				//Draw coords
				if y == 0 {
					buf := fmt.Sprintf("%v", x+1)
					text.Draw(boardCache, buf, monoFontSmall, (mag*x)+offsetPixels+5, (mag*y)+offsetPixels-3, color.White)
				}
				if x == 0 {
					buf := fmt.Sprintf("%2v", y+1)
					text.Draw(boardCache, buf, monoFontSmall, (mag*x)+(offsetPixels/2)+5, (mag*y)+offsetPixels+15, color.White)
				}

				//XY Labels
				text.Draw(boardCache, "X", monoFont, boardPixels/2, 25, color.White)
				text.Draw(boardCache, "Y", monoFont, 5, (boardPixels/2)+65, color.White)
			}
		}

		//Draw right side bg blue
		for x := 0; x < boardSize; x++ {
			for y := 0; y < boardSize; y++ {
				color := ColorBlue
				if (x+y)%2 == 0 {
					color = ColorDarkBlue
				}
				vector.DrawFilledRect(boardCache, float32(mag*x)+(offsetPixels)+boardPixels, float32(mag*y)+offsetPixels, size, size, color, true)
			}
		}

		//Draw board lables
		text.Draw(boardCache, "Audience", monoFont, boardPixels/2, boardPixels+65, color.White)
		text.Draw(boardCache, "Computer", monoFont, boardPixels+boardPixels/2, boardPixels+65, color.White)
		boardCached = true
	}
	screen.DrawImage(boardCache, nil)

	//Draw towers
	for _, item := range gameMap {
		vector.DrawFilledCircle(screen, float32((item.Pos.X+itemOffset)*mag)-(size/1.5), float32((item.Pos.Y+itemOffset)*mag)-(size/1.5), size/2, color.White, true)
	}

}
