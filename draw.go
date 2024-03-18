package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var frameCount uint64

func (g *Game) Draw(screen *ebiten.Image) {
	frameCount++
	drawGameBoard(screen)

	if board.gameover != GAME_RUNNING {
		return
	}

	vector.DrawFilledRect(screen, 0, float32(ScreenHeight)-40, float32(ScreenWidth), 100, ColorSmoke, true)
	if UserMsgDict.VoteState == VOTE_PLAYERS {

		buf := fmt.Sprintf("Vote now: %vx,y -- Votes: %v -- %2.1f remaining%v",
			userSettings.CmdPrefix, UserMsgDict.VoteCount,
			float32(time.Until(UserMsgDict.StartTime.Add(playerRoundTime)).Milliseconds())/1000.0,
			makeEllipsis())

		text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)

	} else if UserMsgDict.VoteState == VOTE_COMPUTER {

		buf := fmt.Sprintf("Computer's turn: %2.1f seconds remaining.",
			float32(time.Until(UserMsgDict.StartTime.Add(cpuRoundTime)).Milliseconds())/1000.0)

		text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)

	} else {
		if !UserMsgDict.GameRunning {

			buf := fmt.Sprintf("No game active%v", makeEllipsis())
			text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)
		}
	}

}

func makeEllipsis() string {
	return strings.Repeat(".", (int(frameCount%120) / 30))
}

/*
func motionSmoothing() {
	if !noSmoothing {
		// Extrapolate position
		startTime = time.Now()
		since := startTime.Sub(lastNetUpdate)
		remaining := FrameSpeedNS - since.Nanoseconds()
		normal = (float64(remaining) / float64(FrameSpeedNS)) - 1

		//Extrapolation limits
		if normal < -1 {
			normal = -1
		} else if normal > 1 {
			normal = 1
		}

		// If there ins't new data yet, extrapolate
		if !dataDirty {
			var smoothPos XY

			//Extrapolated local player position
			smoothPos.X = uint32(float64(oldLocalPlayerPos.X) - ((float64(localPlayerPos.X) - float64(oldLocalPlayerPos.X)) * normal))
			smoothPos.Y = uint32(float64(oldLocalPlayerPos.Y) - ((float64(localPlayerPos.Y) - float64(oldLocalPlayerPos.Y)) * normal))

			//Extrapolated camera position
			sCamPos.X = (uint32(halfScreenX)) + smoothPos.X
			sCamPos.Y = (uint32(halfScreenY)) + smoothPos.Y

			//Extrapolated remote players
			for p, player := range playerList {
				var psmooth XY
				psmooth.X = uint32(float64(player.lastPos.X) - ((float64(player.pos.X) - float64(player.lastPos.X)) * normal))
				psmooth.Y = uint32(float64(player.lastPos.Y) - ((float64(player.pos.Y) - float64(player.lastPos.Y)) * normal))
				playerList[p].spos = XY{X: uint32(psmooth.X), Y: uint32(psmooth.Y)}
			}

			//Extrapolated creatures
			for p, player := range creatureList {
				var psmooth XY
				psmooth.X = uint32(float64(player.lastPos.X) - ((float64(player.pos.X) - float64(player.lastPos.X)) * normal))
				psmooth.Y = uint32(float64(player.lastPos.Y) - ((float64(player.pos.Y) - float64(player.lastPos.Y)) * normal))
				creatureList[p].spos = XY{X: uint32(psmooth.X), Y: uint32(psmooth.Y)}
			}
		}
	} else {
		// Standard mode, just copy data over
		camPos.X = (uint32(halfScreenX)) + localPlayerPos.X
		camPos.Y = (uint32(halfScreenY)) + localPlayerPos.Y
		sCamPos = camPos
	}
}
*/
