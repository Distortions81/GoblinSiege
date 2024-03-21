package main

import (
	"time"
)

func handleMoves() {
	for ServerRunning {
		votes.Lock.Lock()

		//Fast mode for testing quickly, shorten rounds and skip some to get to the action
		if *fastMode {
			if board.moveNum < 20 || (*noTowers && board.moveNum < 40) {
				cpuMoveTime = time.Millisecond * 100
				playerMoveTime = time.Nanosecond
			} else {
				cpuMoveTime = time.Millisecond * 2000
				playerMoveTime = time.Millisecond * 1000
			}
		}

		if votes.VoteState == VOTE_PLAYERS &&
			//Players are voting
			time.Since(votes.StartTime) > playerMoveTime {
			endVote()

		} else if votes.VoteState == VOTE_PLAYERS_DONE {
			//Players are done voting, computer's turn
			votes.VoteState = VOTE_COMPUTER
			votes.CpuTime = time.Now()
			cpuTurn()

		} else if votes.VoteState == VOTE_COMPUTER &&
			//Computer is done, mark new mode
			time.Since(votes.CpuTime) > cpuMoveTime {
			votes.VoteState = VOTE_COMPUTER_DONE
		} else if votes.VoteState == VOTE_COMPUTER_DONE &&
			votes.GameRunning {

			//Computer is done, either start a new vote or skip X rounds
			if board.moveNum%3 == 0 {
				startVote()
			} else {
				votes.VoteState = VOTE_COMPUTER
				votes.CpuTime = time.Now()
				cpuTurn()
			}
		}

		//If a game isn't running, start a new one
		if !votes.GameRunning {
			if time.Since(votes.RoundTime) > time.Second*15 {
				startGame()
			}
		}
		votes.Lock.Unlock()

		//Background wind sound loop
		if !sounds[SND_WIND].player.IsPlaying() {
			sounds[SND_WIND].player.Rewind()
			sounds[SND_WIND].player.SetVolume(sounds[SND_WIND].vol)
			sounds[SND_WIND].player.Play()
		}

		time.Sleep(time.Millisecond * 10)
	}
}

func cpuTurn() {
	board.lock.Lock()

	board.moveNum++

	if board.moveNum >= maxMoves {
		board.gameover = GAME_VICTORY
		playSound(SND_GAMEWON)
		votes.RoundTime = time.Now()
		votes.GameRunning = false
	}

	towerShootArrow()
	spawnGoblins()
	goblinAttack()

	board.lock.Unlock()
}
