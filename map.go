package main

import (
	"fmt"
	"image/color"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
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
)

func getOtype(name string) *oTypeData {
	for o, ot := range oTypes {
		if strings.EqualFold(ot.name, name) {
			return &oTypes[o]
		}
	}
	return nil
}

var oTypes = []oTypeData{
	{name: "Stone Tower", maxHealth: 100, size: xyi{X: 32, Y: 64}, spriteName: "tower1"},
}

type oTypeData struct {
	name       string
	maxHealth  int
	size       xyi
	spriteName string
	spriteImg  *ebiten.Image
}

type objectData struct {
	Pos    xyi
	Health int

	oTypeP *oTypeData
}

var board gameBoardData

const (
	GAME_RUNNING = iota
	GAME_VICTORY
	GAME_DEFEAT
	GAME_DRAW
)

type gameBoardData struct {
	roundNum int
	pmap     map[xyi]*objectData
	emap     map[xyi]*objectData
	lock     sync.Mutex

	gameover int

	bgCache *ebiten.Image
	bgDirty bool
}

func drawGameBoard(screen *ebiten.Image) {

	screen.DrawImage(bgimg, nil)

	if board.bgDirty {
		board.bgCache.Clear()

		for x := 0; x < boardSizeX+enemyBoardX; x++ {
			for y := 0; y < boardSizeY; y++ {

				if x >= boardSizeX {
					if (x+y)%2 == 0 {
						vector.DrawFilledRect(board.bgCache, float32(mag*x)+offPixX, float32(mag*y)+offPixY, size, size, ColorRedC, true)
					}
					continue
				}

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
	if 1 == 1 || UserMsgDict.VoteState == VOTE_PLAYERS {
		screen.DrawImage(board.bgCache, nil)
	}

	//Draw towers
	board.lock.Lock()

	//Works for now, but test if sorted list is faster
	for x := 0; x < boardSizeX; x++ {
		for y := 0; y < boardSizeY; y++ {
			item := board.pmap[xyi{X: x, Y: y}]
			if item == nil {
				continue
			}

			//Draw tower
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(((item.Pos.X+offX)*mag)-item.oTypeP.size.X), float64(((item.Pos.Y+offY)*mag)-item.oTypeP.size.Y))
			screen.DrawImage(item.oTypeP.spriteImg, op)
			//vector.DrawFilledCircle(screen, float32((item.Pos.X+offX)*mag)-(size/2), float32((item.Pos.Y+offY)*mag)-(size/2), size/2, color.White, true)

			//Draw health
			healthBar := (float32(item.Health) / float32(item.oTypeP.maxHealth))
			if healthBar > 0 && healthBar < 1 {
				vector.DrawFilledRect(screen, float32(((item.Pos.X+offX)*mag)-32), float32(((item.Pos.Y+offY)*mag)-64)+1, float32(item.oTypeP.size.X), 6, ColorBlack, false)
				vector.DrawFilledRect(screen, float32(((item.Pos.X+offX)*mag)-31), float32(((item.Pos.Y+offY)*mag)-63)+1, (healthBar*float32(item.oTypeP.size.X) - 1), 4, healthColor(healthBar), false)

			}
		}
	}

	for _, item := range board.emap {
		//Draw tower
		/*
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(((item.Pos.X+offX)*mag)-item.oTypeP.size.X), float64(((item.Pos.Y+offY)*mag)-item.oTypeP.size.Y))
			screen.DrawImage(item.oTypeP.spriteImg, op)
		*/
		vector.DrawFilledCircle(screen, float32((item.Pos.X+offX)*mag)-(size/2), float32((item.Pos.Y+offY)*mag)-(size/2), size/2, ColorRed, true)

		//Draw health
		/*
			healthBar := (float32(item.Health) / float32(item.oTypeP.maxHealth))
			if healthBar > 0 && healthBar < 1 {
				vector.DrawFilledRect(screen, float32(((item.Pos.X+offX)*mag)-32), float32(((item.Pos.Y+offY)*mag)-64)+1, float32(item.oTypeP.size.X), 6, ColorBlack, false)
				vector.DrawFilledRect(screen, float32(((item.Pos.X+offX)*mag)-31), float32(((item.Pos.Y+offY)*mag)-63)+1, (healthBar*float32(item.oTypeP.size.X) - 1), 4, healthColor(healthBar), false)

			}
		*/
	}

	if board.gameover == GAME_DEFEAT {
		vector.DrawFilledRect(screen, 0, float32(ScreenHeight)-40, float32(ScreenWidth), 100, ColorSmoke, true)
		buf := fmt.Sprintf("Game over: The audience was defeated on round %v!", board.roundNum)
		text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)
	}

	board.lock.Unlock()

}

func healthColor(input float32) color.NRGBA {
	var healthColor color.NRGBA = color.NRGBA{R: 255, G: 255, B: 255, A: 0}
	health := input * 100

	if health < 100 && health > 0 {
		healthColor.A = 255
		healthColor.B = 0

		r := int(float32(100-(health)) * 5)
		if r > 255 {
			r = 255
		}
		healthColor.R = uint8(r)

		g := int(float32(health) * 4)
		if g > 255 {
			g = 255
		}
		healthColor.G = uint8(g)

	}

	return healthColor
}
