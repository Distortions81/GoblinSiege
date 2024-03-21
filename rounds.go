package main

import (
	"time"
)

func handleMoves() {
	for ServerRunning {
		votes.Lock.Lock()

		if *fastMode {
			if board.moveNum < 20 {
				cpuMoveTime = time.Millisecond * 100
				playerMoveTime = time.Nanosecond
			} else {
				cpuMoveTime = time.Millisecond * 2000
				playerMoveTime = time.Millisecond * 1000
			}
		}

		if votes.VoteState == VOTE_PLAYERS &&
			time.Since(votes.StartTime) > playerMoveTime {
			endVote()

		} else if votes.VoteState == VOTE_PLAYERS_DONE {
			votes.VoteState = VOTE_COMPUTER
			votes.CpuTime = time.Now()
			cpuTurn()

		} else if votes.VoteState == VOTE_COMPUTER &&
			time.Since(votes.CpuTime) > cpuMoveTime {

			votes.VoteState = VOTE_COMPUTER_DONE

		} else if votes.VoteState == VOTE_COMPUTER_DONE &&
			votes.GameRunning {
			if board.moveNum%3 == 0 {
				startVote()
			} else {
				votes.VoteState = VOTE_COMPUTER
				votes.CpuTime = time.Now()
				cpuTurn()
			}
		}

		if !votes.GameRunning {
			if time.Since(votes.RoundTime) > time.Second*5 {
				startGame()
			}
		}

		votes.Lock.Unlock()
		time.Sleep(time.Millisecond * 10)

		//Bg sound
		if !sounds[SND_WIND].player.IsPlaying() {
			sounds[SND_WIND].player.Rewind()
			sounds[SND_WIND].player.SetVolume(sounds[SND_WIND].vol)
			sounds[SND_WIND].player.Play()
		}
	}
}

func cpuTurn() {
	board.lock.Lock()

	board.moveNum++

	if board.moveNum >= maxMoves {
		board.gameover = GAME_VICTORY
		endGame()
	}

	towerShootArrow()
	spawnGoblins()
	goblinAttack()

	board.lock.Unlock()
}
