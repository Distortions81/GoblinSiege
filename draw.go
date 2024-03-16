package main

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hako/durafmt"
)

var frameCount uint64

func (g *Game) Draw(screen *ebiten.Image) {
	frameCount++
	drawGameBoard(screen)

	if UserMsgDict.VoteState == VOTE_PLAYERS {
		vector.DrawFilledRect(screen, 0, float32(ScreenHeight)-50, 1280, 100, ColorSmoke, true)
		buf := fmt.Sprintf("Vote now: %vx,y -- Votes: %v -- %v remaining%v",
			userSettings.CmdPrefix, UserMsgDict.VoteCount,
			durafmt.Parse(time.Until(UserMsgDict.StartTime.Add(playerRoundTime)).Round(time.Second)).LimitFirstN(1),
			makeEllipsis())

		text.Draw(screen, buf, monoFont, 10, ScreenHeight-50+35, color.White)

	} else if UserMsgDict.VoteState == VOTE_COMPUTER {
		vector.DrawFilledRect(screen, 0, float32(ScreenHeight)-50, 1280, 100, ColorSmoke, true)
		buf := fmt.Sprintf("Computer's turn: %v remaining%v",
			durafmt.Parse(time.Until(UserMsgDict.StartTime.Add(cpuRoundTime)).Round(time.Second)).LimitFirstN(1),
			makeEllipsis())

		text.Draw(screen, buf, monoFont, 10, ScreenHeight-15, color.White)

	} else {
		if !UserMsgDict.GameRunning {
			text.Draw(screen, "No game active.", monoFont, 10, ScreenHeight-7, color.White)
		}
	}

}

func makeEllipsis() string {
	return strings.Repeat(".", (int(frameCount%120) / 30))
}
