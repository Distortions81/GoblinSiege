package main

import (
	"fmt"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

type objectData struct {
	Pos    xyi
	OldPos xyi

	Health int
	dead   bool

	sheetP *spriteSheetData
}

var board gameBoardData

const (
	GAME_RUNNING = iota
	GAME_VICTORY
	GAME_DEFEAT
	GAME_DRAW
)

type arrowData struct {
	tower  xyi
	target xyi

	shot   time.Time
	missed bool
}

type gameBoardData struct {
	moveNum  int
	playMap  map[xyi]*objectData
	enemyMap map[xyi]*objectData
	lock     sync.Mutex

	arrowsShot []arrowData

	gameover int

	bgCache *ebiten.Image
	bgDirty bool
}

func drawGameBoard(screen *ebiten.Image) {

	screen.DrawImage(bgimg, nil)

	//Draw checkerboard if dirty
	if board.bgDirty {
		board.bgCache.Clear()

		for x := 0; x < boardSizeX; x++ {
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
	//Draw checkerboard cache if voting
	if votes.VoteState == VOTE_PLAYERS {
		screen.DrawImage(board.bgCache, nil)
	}

	//Draw board
	board.lock.Lock()
	defer board.lock.Unlock()

	//Draw arrows
	numArrows := len(board.arrowsShot) - 1
	startTime := time.Now()
	for x := numArrows; x >= 0; x-- {
		arrow := board.arrowsShot[x]

		//Tween animation, make sprite face direction of travel
		since := startTime.Sub(arrow.shot)
		distance := Distance(arrow.tower, arrow.target)
		const ratio = 30
		remaining := (distance * float64(cpuMoveTime.Nanoseconds()/ratio)) - float64(since.Nanoseconds())
		normal := (float64(remaining)/(distance*float64(cpuMoveTime.Nanoseconds()/ratio)) - 1.0)

		//Extrapolation limits
		if normal < -1 {
			normal = -1
		} else if normal > 1 {
			normal = 1
		}

		sX := (float64(arrow.tower.X) - ((float64(arrow.target.X) - float64(arrow.tower.X)) * normal))
		sY := (float64(arrow.tower.Y) - ((float64(arrow.target.Y) - float64(arrow.tower.Y)) * normal))

		//Hide arrows that didn't miss once at target
		if sX == float64(arrow.target.X) && sY == float64(arrow.target.Y) {
			if !arrow.missed {
				//Delete it
				board.arrowsShot = append(board.arrowsShot[:x], board.arrowsShot[x+1:]...)
				continue
			}
		}

		towerPos := geom.Coord{float64(arrow.tower.X), float64(arrow.tower.Y), 0}
		targetPos := geom.Coord{float64(arrow.target.X), float64(arrow.target.Y), 0}
		angle := xy.Angle(towerPos, targetPos)

		//Draw arrow
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate(angle)
		op.GeoM.Translate(((sX+float64(offX))*float64(mag))-float64(obj_arrow.frameSize.X)-16,
			((sY+float64(offY))*float64(mag))-float64(obj_arrow.frameSize.Y)-16)

		screen.DrawImage(obj_arrow.img, op)
	}

	//Draw goblin
	for _, item := range board.enemyMap {
		//Tween animation

		since := startTime.Sub(votes.CpuTime)
		remaining := (float64(cpuMoveTime.Nanoseconds())) - float64(since.Nanoseconds())
		normal := (float64(remaining)/(float64(cpuMoveTime.Nanoseconds())) - 1.0)

		//Extrapolation limits
		if normal < -1 {
			normal = -1
		} else if normal > 1 {
			normal = 1
		}

		sX := (float64(item.OldPos.X) - ((float64(item.Pos.X) - float64(item.OldPos.X)) * normal))
		sY := (float64(item.OldPos.Y) - ((float64(item.Pos.Y) - float64(item.OldPos.Y)) * normal))

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(((sX+float64(offX))*float64(mag))-float64(obj_arrow.frameSize.X),
			((sY+float64(offY))*float64(mag))-float64(obj_arrow.frameSize.Y))

		if item.dead {
			screen.DrawImage(item.sheetP.anims[ANI_DIE].img[int(sX*16)%4], op)
		} else {
			screen.DrawImage(item.sheetP.anims[ANI_RUN].img[int(sX*16)%4], op)
			healthBar := (float32(item.Health) / float32(item.sheetP.health))

			if healthBar > 0 && healthBar < 1 {
				vector.DrawFilledRect(screen, float32(((sX+offX)*mag)-32), float32(((sY+offY)*mag)-32)+1, float32(item.sheetP.frameSize.X), 6, ColorBlack, false)
				vector.DrawFilledRect(screen, float32(((sX+offX)*mag)-31), float32(((sY+offY)*mag)-31)+1, (healthBar*float32(item.sheetP.frameSize.X) - 1), 4, healthColor(healthBar), false)
			}
		}
	}

	//Draw towers
	for x := 0; x <= boardSizeX; x++ {
		for y := 0; y <= boardSizeY; y++ {
			item := board.playMap[xyi{X: x, Y: y}]
			if item == nil {
				continue
			}

			//Draw tower
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(((item.Pos.X+offX)*mag)-item.sheetP.frameSize.X), float64(((item.Pos.Y+offY)*mag)-item.sheetP.frameSize.Y))
			if item.dead {
				screen.DrawImage(item.sheetP.img, op)
			} else {
				screen.DrawImage(item.sheetP.img, op)

				//Draw health
				healthBar := (float32(item.Health) / float32(item.sheetP.health))
				if healthBar > 0 && healthBar < 1 {
					vector.DrawFilledRect(screen, float32(((item.Pos.X+offX)*mag)-32), float32(((item.Pos.Y+offY)*mag)-64)+1, float32(item.sheetP.frameSize.X), 6, ColorBlack, false)
					vector.DrawFilledRect(screen, float32(((item.Pos.X+offX)*mag)-31), float32(((item.Pos.Y+offY)*mag)-63)+1, (healthBar*float32(item.sheetP.frameSize.X) - 1), 4, healthColor(healthBar), false)

				}
			}

		}
	}

	if board.gameover == GAME_DEFEAT {
		vector.DrawFilledRect(screen, 0, float32(ScreenHeight)-40, float32(ScreenWidth), 100, ColorSmoke, true)
		buf := fmt.Sprintf("Game over: The audience was defeated on move %v!", board.moveNum)
		text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)
	} else if board.gameover == GAME_VICTORY {
		vector.DrawFilledRect(screen, 0, float32(ScreenHeight)-40, float32(ScreenWidth), 100, ColorSmoke, true)
		buf := fmt.Sprintf("Game over: The audience has won, survived %v move!", board.moveNum)
		text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)
	}

	buf := fmt.Sprintf("Move: %v/%v!", board.moveNum, maxMoves)
	text.Draw(screen, buf, monoFont, ScreenWidth-210, 25, color.Black)
}
