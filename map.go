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
	boardSizeX   = 28
	boardSizeY   = 19
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

	screen.DrawImage(bgimg, nil)

	if board.bgDirty {
		board.bgCache.Clear()

		//Draw left side bg red

		for x := 0; x < boardSizeX; x++ {
			for y := 0; y < boardSizeY; y++ {

				if (x+y)%2 == 0 {
					vector.DrawFilledRect(board.bgCache, float32(mag*x)+offPixX, float32(mag*y)+offPixY, size, size, ColorGreenC, true)
				}

				//Draw coords

				if x == 0 {
					buf := fmt.Sprintf("%2v", y+1)
					text.Draw(board.bgCache, buf, monoFontSmall, offPixX-(mag/2), (mag*y)+offPixY+20, color.Black)
				}
				if y == 0 {
					buf := fmt.Sprintf("%v", x+1)
					text.Draw(board.bgCache, buf, monoFontSmall, (mag*x)+offPixX+8, (mag*y)+offPixY-2, color.Black)
				}

				//XY Labels
				text.Draw(board.bgCache, "X", monoFont, offPixX+(boardPixelsX/2), 20, color.Black)
				text.Draw(board.bgCache, "Y", monoFont, offPixX-(mag), offPixY+(boardPixelsY/2), color.Black)
			}
		}

		board.bgDirty = false
	}
	if UserMsgDict.VoteState == VOTE_PLAYERS {
		screen.DrawImage(board.bgCache, nil)
	}

	//Draw towers
	board.lock.Lock()
	for _, item := range board.bmap {
		vector.DrawFilledCircle(screen, float32((item.Pos.X+offX)*mag)-(size/1.5), float32((item.Pos.Y+offY)*mag)-(size/1.5), size/2, color.White, true)
	}
	board.lock.Unlock()

}
