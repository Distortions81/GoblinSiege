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
	size         = 32
	mag          = size
	boardSizeX   = 27
	boardSizeY   = 17
	offX         = 5
	offY         = 1
	offPixX      = size * offX
	offPixY      = size * offY
	boardPixelsX = ((boardSizeX) * mag)
	boardPixelsY = ((boardSizeY) * mag)
)

type objectData struct {
	Pos  xyi
	Type int
}

var board gameBoardData

type gameBoardData struct {
	roundNum int
	bmap     map[xyi]*objectData
	lock     sync.Mutex

	bgCache *ebiten.Image
	bgDirty bool
}

func drawGameBoard(screen *ebiten.Image) {

	if board.bgDirty {
		board.bgCache.Clear()

		board.bgCache.DrawImage(bgimg, nil)

		//Draw left side bg red
		for x := 0; x < boardSizeX; x++ {
			for y := 0; y < boardSizeY; y++ {
				var tColor color.Color

				tColor = ColorDarkGreen
				if (x+y)%2 == 0 {
					tColor = ColorGreen
				}

				vector.DrawFilledRect(board.bgCache, float32(mag*x)+offPixX, float32(mag*y)+offPixY, size, size, tColor, true)

				//Draw coords
				if y == 0 {
					buf := fmt.Sprintf("%v", x+1)
					text.Draw(board.bgCache, buf, monoFontSmall, (mag*x)+offPixX+5, (mag*y)+offPixY-3, color.White)
				}
				if x == 0 {
					buf := fmt.Sprintf("%2v", y+1)
					text.Draw(board.bgCache, buf, monoFontSmall, (mag*x)+(offPixX/2)+5, (mag*y)+offPixY+15, color.White)
				}

				//XY Labels
				text.Draw(board.bgCache, "X", monoFont, boardPixelsX/2, 25, color.White)
				text.Draw(board.bgCache, "Y", monoFont, 5, (boardPixelsY/2)+65, color.White)
			}
		}

		board.bgDirty = false
	}
	screen.DrawImage(board.bgCache, nil)

	//Draw towers
	board.lock.Lock()
	for _, item := range board.bmap {
		vector.DrawFilledCircle(screen, float32((item.Pos.X+offX)*mag)-(size/1.5), float32((item.Pos.Y+offY)*mag)-(size/1.5), size/2, color.White, true)
	}
	board.lock.Unlock()

}
