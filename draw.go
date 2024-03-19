package main

import (
	"fmt"
	"image/color"
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
	if votes.VoteState == VOTE_PLAYERS {
		till := float32(time.Until(votes.StartTime.Add(playerMoveTime)).Milliseconds()) / 1000.0
		if till > 0 {
			buf := fmt.Sprintf("Your turn!!! Vote: %vx,y -- Votes: %v -- %2.1fs remaining%v",
				userSettings.CmdPrefix, votes.VoteCount,
				till,
				makeEllipsis())

			text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)
		}

	} else if votes.VoteState == VOTE_COMPUTER || votes.VoteState == VOTE_COMPUTER_DONE {

		till := float32(time.Until(votes.StartTime.Add(cpuMoveTime*3)).Milliseconds()) / 1000.0
		if till > 0 {
			buf := fmt.Sprintf("ORC'S TURN. Vote: %vx,y -- Votes: %v -- %2.1fs remaining%v",
				userSettings.CmdPrefix, votes.VoteCount,
				till,
				makeEllipsis())

			text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)
		}
	} else {
		if !votes.GameRunning {

			buf := fmt.Sprintf("No game active%v", makeEllipsis())
			text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)
		}
	}

}
