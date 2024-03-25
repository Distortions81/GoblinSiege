package main

import (
	"fmt"
	"image/color"
	"sync/atomic"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	//Doesn't really need to be atomic, but this keeps it our of the -race logs
	aniCount atomic.Uint64
)

func (g *Game) Draw(screen *ebiten.Image) {

	if gameMode == MODE_SPLASH {
		screen.DrawImage(splash, nil)

		for _, button := range splashButtons {
			if time.Since(button.clicked) < flashSpeed {
				vector.DrawFilledRect(screen, button.pos.X, button.pos.Y, button.size.X, button.size.Y, ColorRedC, false)
			} else if button.hover {
				vector.DrawFilledRect(screen, button.pos.X, button.pos.Y, button.size.X, button.size.Y, ColorWhiteSmoke, false)
			}
		}

	} else if gameMode == MODE_PLAY_TWITCH ||
		gameMode == MODE_PLAY_SINGLE {

		gameLock.Lock()
		defer gameLock.Unlock()

		// If there isn't a game running, don't render game board
		// Render to an image and fade out at game end
		if board.gameover != GAME_RUNNING {
			if !board.useFreeze {
				//Draw actual game board
				drawGameBoard(board.fFrame)
				board.useFreeze = true
			}
			screen.DrawImage(bgimg, nil)

			op := &ebiten.DrawImageOptions{}
			shotAgo := time.Since(votes.RoundTime)
			pa := 1.0 - float32(shotAgo.Seconds()/gameOverFadeSec)
			if pa > 0 {
				op.ColorScale.ScaleAlpha(pa)
				screen.DrawImage(board.fFrame, op)
			}

			//Handle game ending conditions
			if board.gameover == GAME_DEFEAT {
				vector.DrawFilledRect(screen, 0, float32(defaultWindowHeight)-40, float32(defaultWindowWidth), 100, ColorSmoke, true)
				buf := fmt.Sprintf("GAME OVER: The audience was defeated! Enemy won on move %v!", board.moveNum)
				text.Draw(screen, buf, monoFont, 10, defaultWindowHeight-15, color.White)
			} else if board.gameover == GAME_VICTORY {
				vector.DrawFilledRect(screen, 0, float32(defaultWindowHeight)-40, float32(defaultWindowWidth), 100, ColorSmoke, true)
				buf := fmt.Sprintf("GAME OVER: The audience has won! Survived %v moves!", board.moveNum)
				text.Draw(screen, buf, monoFont, 10, defaultWindowHeight-15, color.White)
			}
			return
		} else {
			drawGameBoard(screen)
		}

		//Show recent votes
		numV := len(newVoteNotice) - 1
		lBuf := "Recent votes:\n"
		for v := numV; v >= 0 && v > (numV-4); v-- {
			voter := newVoteNotice[v]
			if time.Since(voter.time) > time.Second*30 {
				newVoteNotice = append(newVoteNotice[:v], newVoteNotice[v+1:]...)
			} else {
				buf := fmt.Sprintf("%v: %v,%v for move %v\n", voter.sender, voter.pos.X, voter.pos.Y, voter.playerMove)
				lBuf = lBuf + buf
			}
		}
		text.Draw(screen, lBuf, monoFontSmall, (boardSizeX+offX)*mag+50, 60, color.Black)

		if votes.VoteState == VOTE_PLAYERS {
			//Draw player vote
			vector.DrawFilledRect(screen, 0, float32(defaultWindowHeight)-40, float32(defaultWindowWidth), 100, ColorRed, true)

			till := float32(time.Until(votes.StartTime.Add(playerMoveTime)).Milliseconds()) / 1000.0
			if till > 0 {
				buf := fmt.Sprintf("Your turn! -- Place or upgrade tower: %vx,y (%v votes)%v",
					userSettings.CmdPrefix, votes.VoteCount, makeEllipsis())

				text.Draw(screen, buf, monoFont, 10, defaultWindowHeight-15, color.Black)
			}

		} else if votes.VoteState == VOTE_COMPUTER || votes.VoteState == VOTE_COMPUTER_DONE {
			//Draw enemy turn and background voting
			vector.DrawFilledRect(screen, 0, float32(defaultWindowHeight)-40, float32(defaultWindowWidth), 100, ColorSmoke, true)

			till := float32(time.Until(votes.StartTime.Add(cpuMoveTime*3)).Milliseconds()) / 1000.0
			if till > 0 {
				buf := fmt.Sprintf("Enemy's turn! -- Vote to place or upgrade tower: %vx,y (%v votes)%v",
					userSettings.CmdPrefix, votes.VoteCount, makeEllipsis())

				text.Draw(screen, buf, monoFont, 10, defaultWindowHeight-15, color.White)
			}
		} else {
			//No game active
			if !votes.GameRunning {
				buf := "No game active."
				text.Draw(screen, buf, monoFont, 10, defaultWindowHeight-15, color.White)
			}
		}

		if *debugMode {
			buf := fmt.Sprintf("%2.2f fps,%v towers, %v goblin, v%v, ",
				ebiten.ActualFPS(),
				len(board.towerMap),
				len(board.goblinMap),
				version)
			text.Draw(screen, buf, monoFont, 10, 24, ColorBlack)
		}
	}
}
