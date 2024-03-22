package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

const (
	OTYPE_TOWER = iota
	OTYPE_VWALL
	OTYPE_MAX
)

type objectData struct {
	worldObjType int
	pos          xyi
	prevPos      xyi

	health int
	dead   bool
	diedAt time.Time

	aniOffset uint64

	attacking    bool
	lastAttacked time.Time
	building     int
	upgrade      int

	//Item spritesheet data
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

	arrowsShot  []arrowData
	wallDmgTime time.Time

	gameover int
	bgCache  *ebiten.Image
	bgDirty  bool
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

	//Draw wall covering
	for y := 0; y <= boardSizeY; y++ {
		item := board.playMap[xyi{X: -1, Y: y}]
		if item == nil {
			continue
		}
		if item.worldObjType == OTYPE_VWALL {

			op := &ebiten.DrawImageOptions{}
			if time.Since(item.lastAttacked) < flashSpeed {
				op.ColorScale.Scale(2, 0.5, 0.5, 1)
			} else if !item.dead {
				continue
			}

			op.GeoM.Translate(float64(((item.pos.X+offX)*mag)-item.sheetP.frameSize.X), float64(((item.pos.Y+offY)*mag)-item.sheetP.frameSize.Y))
			screen.DrawImage(item.sheetP.img, op)
		}
	}
	//Draw arrows
	numArrows := len(board.arrowsShot) - 1
	startTime := time.Now()
	for x := numArrows; x >= 0; x-- {
		arrow := board.arrowsShot[x]

		//Tween animation, make sprite face direction of travel
		since := startTime.Sub(arrow.shot)
		distance := Distance(arrow.tower, arrow.target)
		remaining := (distance * float64(cpuMoveTime.Nanoseconds()/arrowSpeed)) - float64(since.Nanoseconds())
		normal := (float64(remaining)/(distance*float64(cpuMoveTime.Nanoseconds()/arrowSpeed)) - 1.0)

		//Extrapolation limits
		if normal < -1 {
			normal = -1
		} else if normal > 1 {
			normal = 1
		}

		//Tweened position
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

		//Tweening begin and ending points, convert to geom.Coord for the xy library
		towerPos := geom.Coord{float64(arrow.tower.X), float64(arrow.tower.Y), 0}
		targetPos := geom.Coord{float64(arrow.target.X), float64(arrow.target.Y), 0}
		angle := xy.Angle(towerPos, targetPos)

		//Draw arrow
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate(angle)
		op.GeoM.Translate(((sX+float64(offX))*float64(mag))-float64(obj_arrow.frameSize.X)-16,
			((sY+float64(offY))*float64(mag))-float64(obj_arrow.frameSize.Y)-16)

		var pa float32
		shotAgo := time.Since(arrow.shot)
		pa = 1.0 - float32(shotAgo.Seconds()/arrowFadeSec)
		if pa < 0 {
			//Delete it
			board.arrowsShot = append(board.arrowsShot[:x], board.arrowsShot[x+1:]...)
			continue
		}

		op.ColorScale.ScaleAlpha(pa)
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

		//Tweened position
		sX := (float64(item.prevPos.X) - ((float64(item.pos.X) - float64(item.prevPos.X)) * normal))
		sY := (float64(item.prevPos.Y) - ((float64(item.pos.Y) - float64(item.prevPos.Y)) * normal))

		op := &ebiten.DrawImageOptions{}
		//Horizontal mirroring for sprites that are marked mirror
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(((sX + float64(offX)) * float64(mag)),
			((sY+float64(offY))*float64(mag))-float64(obj_goblinBarb.frameSize.Y))

		if item.pos.X > 31 {
			//Swim animation for the water
			screen.DrawImage(item.sheetP.anims[ANI_SWIM].img[int(float64(item.aniOffset)+sX*16)%4], op)

		} else if item.dead && time.Since(item.diedAt) > deathDelay {
			//Goblin died
			deadAni := 0
			if time.Since(item.diedAt) > (deathDelay) {
				deadAni = 1
			}
			var pa float32
			deadFor := time.Since(item.diedAt).Seconds()
			pa = 1.0 - float32(deadFor/bodyFadeSec)
			if pa < 0 {
				continue
			}
			op.ColorScale.ScaleAlpha(pa)
			screen.DrawImage(item.sheetP.anims[ANI_DIE].img[deadAni], op)

		} else if item.attacking {
			//Goblin attacking tower
			attackFrame := time.Since(votes.CpuTime) / attackDelay
			if attackFrame > 3 {
				attackFrame = 3
			}
			screen.DrawImage(item.sheetP.anims[ANI_ATTACK].img[attackFrame%4], op)
		} else {
			///Draw running
			screen.DrawImage(item.sheetP.anims[ANI_RUN].img[int(float64(item.aniOffset)+(sX+offX)*16)%4], op)
		}

		//Show health bar
		healthBar := (float32(item.health) / float32(item.sheetP.health))
		if !item.dead && healthBar > 0 && healthBar < 1 {
			vector.DrawFilledRect(screen, float32(((sX+offX)*mag)-32), float32(((sY+offY)*mag)-32)+1, float32(item.sheetP.frameSize.X), 4, ColorSmoke, false)
			vector.DrawFilledRect(screen, float32(((sX+offX)*mag)-31), float32(((sY+offY)*mag)-31)+1, (healthBar*float32(item.sheetP.frameSize.X) - 1), 2, healthColor(healthBar), false)
		}
	}

	//Draw towers
	for x := -1; x <= boardSizeX; x++ {
		for y := 0; y <= boardSizeY; y++ {
			item := board.playMap[xyi{X: x, Y: y}]
			if item == nil {
				continue
			}

			if item.worldObjType == OTYPE_TOWER {
				//Draw tower
				op := &ebiten.DrawImageOptions{}
				if time.Since(item.lastAttacked) < flashSpeed {
					op.ColorScale.Scale(2, 0.5, 0.5, 1)
				}
				op.GeoM.Translate(float64(((item.pos.X+offX)*mag)-item.sheetP.frameSize.X), float64(((item.pos.Y+offY)*mag)-item.sheetP.frameSize.Y))
				if item.dead {
					//Broken tower
					screen.DrawImage(item.sheetP.anims[ANI_FADE].img[(aniCount+item.aniOffset)%3], op)
				} else {
					//Draw tower being built, otherwise animate fully built one
					if item.building < 2 {
						screen.DrawImage(item.sheetP.anims[ANI_RUN].img[item.building%3], op)
					} else {
						screen.DrawImage(item.sheetP.anims[ANI_IDLE].img[(aniCount+item.aniOffset)%3], op)
					}
				}
			}

			//Draw health
			healthBar := (float32(item.health) / float32(item.sheetP.health))
			if healthBar > 0 && healthBar < 1 {
				vector.DrawFilledRect(screen, float32(((item.pos.X+offX)*mag)-32), float32(((item.pos.Y+offY)*mag)-64)+1, float32(item.sheetP.frameSize.X), 4, ColorSmoke, false)
				vector.DrawFilledRect(screen, float32(((item.pos.X+offX)*mag)-31), float32(((item.pos.Y+offY)*mag)-63)+1, (healthBar*float32(item.sheetP.frameSize.X) - 1), 2, healthColor(healthBar), false)
			}
		}
	}

	//Show the current move number in the corner
	buf := fmt.Sprintf("Move: %v/%v!", board.moveNum, maxMoves)
	text.Draw(screen, buf, monoFont, defaultWindowWidth-210, 25, color.Black)
}
