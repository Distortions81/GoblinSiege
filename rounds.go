package main

import "time"

func handleRounds() {
	for ServerRunning {
		UserMsgDict.Lock.Lock()

		if UserMsgDict.VoteState == VOTE_PLAYERS &&
			time.Since(UserMsgDict.StartTime) > playerRoundTime {
			endVote()

		} else if UserMsgDict.VoteState == VOTE_PLAYERS_DONE {
			UserMsgDict.VoteState = VOTE_COMPUTER
			UserMsgDict.StartTime = time.Now()
			cpuTurn()

		} else if UserMsgDict.VoteState == VOTE_COMPUTER &&
			time.Since(UserMsgDict.StartTime) > cpuRoundTime {

			UserMsgDict.VoteState = VOTE_COMPUTER_DONE

		} else if UserMsgDict.VoteState == VOTE_COMPUTER_DONE &&
			UserMsgDict.GameRunning {
			startVote()
		}

		UserMsgDict.Lock.Unlock()
		time.Sleep(time.Millisecond * 100)
	}
}

func cpuTurn() {
	board.lock.Lock()
	defer board.lock.Unlock()

	board.roundNum++

}
