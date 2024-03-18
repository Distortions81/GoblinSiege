package main

import (
	"log"
	"math/rand"
	"time"
)

func addTower() {
	tower1 := getOtype("Stone Tower")

	if UserMsgDict.VoteCount > 0 &&
		UserMsgDict.Result.X > 0 &&
		UserMsgDict.Result.Y > 0 &&
		UserMsgDict.Result.X <= boardSizeX &&
		UserMsgDict.Result.Y <= boardSizeY {

		tpos := UserMsgDict.Result
		if board.emap[tpos] == nil && board.pmap[tpos] == nil {
			board.pmap[tpos] = &objectData{Pos: tpos, oTypeP: tower1, Health: tower1.maxHealth}
		} else {
			log.Println("COLLISION!")
		}
	} else {

		log.Println("Not enough votes, picking random.")
		//Invalid or not enough votes, pick a pos at random
		tpos := xyi{X: rand.Intn(boardSizeX-1) + 1, Y: rand.Intn(boardSizeY-1) + 1}
		if board.emap[tpos] == nil && board.pmap[tpos] == nil {
			board.pmap[tpos] = &objectData{Pos: tpos, oTypeP: tower1, Health: tower1.maxHealth}
		}
	}

}

func towerShootArrow() {
	curTime := time.Now()

	for _, item := range board.pmap {
		if item.dead {
			continue
		}
		for _, enemy := range board.emap {
			if enemy.dead {
				continue
			}
			//If enemy within range
			if Distance(item.Pos, enemy.Pos) < 6 {

				if rand.Intn(2) != 0 {
					arrow := arrowData{tower: item.Pos, target: enemy.Pos, missed: true, shot: curTime}
					board.arrowsShot = append(board.arrowsShot, arrow)
					break
				}
				arrow := arrowData{tower: item.Pos, target: enemy.Pos, missed: false, shot: curTime}
				board.arrowsShot = append(board.arrowsShot, arrow)

				dmgAmt := 5 + rand.Intn(20)
				enemy.Health -= dmgAmt

				if enemy.Health <= 0 {
					board.emap[enemy.Pos].dead = true
					//For tweening
					board.emap[enemy.Pos].OldPos = board.emap[enemy.Pos].Pos
				}
				break
			}
		}
	}
}
